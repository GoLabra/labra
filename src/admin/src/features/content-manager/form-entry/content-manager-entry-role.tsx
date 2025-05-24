import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { forwardRef, useCallback, useEffect, useImperativeHandle, useMemo } from "react";
import { ContentManagerEntryDialogContentProps } from "../content-manager-entry-form";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { BasicAuditTrail } from "@/shared/components/basic-audit-trail";
import { Id } from "@/shared/components/id";
import { zodResolver } from "@hookform/resolvers/zod";
import { DialogContent, Stack, Button, Typography, Checkbox, Box } from "@mui/material";
import { useForm } from "react-hook-form";
import { useContentManagerFormSchema } from "../use-content-manager-form-schema";
import { Form } from "@/core-features/dynamic-form2/dynamic-form";
import { z } from "zod";
import { TextShortFormField } from "@/core-features/dynamic-form/form-fields/TextShortField";
import { BooleanFormComponent, BooleanFormField } from "@/core-features/dynamic-form/form-fields/BooleanField";
import { SectionTitle } from "@/shared/components/section-title";
import { gql, useLazyQuery, useQuery } from "@apollo/client";
import { Permission, Role } from "@/lib/apollo/graphql";
import { EdgeStatus } from "@/lib/utils/edge-status";
import { groupByMap } from "@/lib/utils/array";
import { PiMathOperations } from "react-icons/pi";
import { GenericEvent } from "@/lib/utils/event";


const GET_ROLE_PERMISSION_QUERY = gql`query getRolePermissionQuery($where: RoleWhereInput) {
	roles(where: $where)  {
		name
		permissions {
    		entity
    		operation
    	}
  	}
}`


type Operation = {
	operation: string;
	status: EdgeStatus;
}

type EntityPermissions = {
	entityName: string;
	operation: Operation[];
}

type OperationDef = {
	name: string;
	caption: string;
}

// PERMISSION SECTION
interface EntityPermissionViewProps {
	entityCaption: string;
	operationDefs: OperationDef[];
	operations: string[];
	valuesChanged?: (operations: string[]) => void;
}
const EntityPermissionView = (props: EntityPermissionViewProps) => {

	const onChange = useCallback((event: GenericEvent<boolean>, operationName: string) => {
		if(event.target.value) {
			props.valuesChanged?.([...props.operations, operationName]);
			return;
		}

		props.valuesChanged?.(props.operations.filter(i => i !== operationName));
	}, [props.valuesChanged]);

	return (
		<Stack gap={1.5} p={1} my={2} borderRadius={0}>
			<Typography >
				{props.entityCaption}
			</Typography>

			<Stack direction="row" gap={1} flexWrap="wrap">
				{props.operationDefs.map((i: OperationDef) => {
					return <BooleanFormComponent 
								key={i.name} 
								name={`${props.entityCaption}-${i.name}`} 
								label={i.caption}
								value={props.operations.includes(i.name)}
								onChange={(event) => onChange(event, i.name)}
								 />
				})}
			</Stack>
		</Stack>
	)
}

export const permissionSchema = z.object({
  caption4: z.string().nonempty('Caption 4 is required'),
  caption5: z.string().nonempty('Caption 5 is required'),
  caption6: z.string().nonempty('Caption 6 is required'),
  required: z.coerce.boolean().optional()
});

interface PermissionSectionProps {
	values: Record<string, Operation[]>;
	valuesChanged?: (entityName: string,operations: string[]) => void;
	
}
const PermissionSection = (props: PermissionSectionProps) => {

	const { entities } = useEntities();

	const operationDefs = [{
		caption: 'Create',
		name: 'Create'
	},
	{
		caption: 'Read',
		name: 'Read'
	},
	{
		caption: 'Update',
		name: 'Update'
	},
	{
		caption: 'Delete',
		name: 'Delete'
	}];

	return (
		<Stack gap={1.5}>
			
			<SectionTitle title="Permissions" />

			<Box sx={{
				'> *:nth-child(even)': {
					backgroundColor: '#99999920'
				}
			}}>

				{entities.map((entity) => {
					return <EntityPermissionView

						key={entity.name} 
						entityCaption={entity.caption} 
						operationDefs={operationDefs} 
						operations={props.values[entity.name]?.map(i => i.operation) ?? []} 
						valuesChanged={(values) => {
							props.valuesChanged?.(entity.name, values);
						}} />
				})}
			</Box>


		</Stack>
	)
}
// END OF PERMISSION SECTION



export const schema = z.object({
	caption: z.string().nonempty('Caption is required'),
	required: z.coerce.boolean().optional()
}).merge(permissionSchema);

export const ContentManagerEntryRole = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	const myDialogContext = useMyDialogContext();
	const fullEntity = useFullEntity({ entityName: props.entityName });
	const formSchema = useContentManagerFormSchema(fullEntity);

	const permissionRequest = useQuery<{ roles: Role[] }>(GET_ROLE_PERMISSION_QUERY, {
		variables: {
			where: {
				name: 'Ham'
			}
		},
		skip: myDialogContext.openMode === FormOpenMode.New
	});

	const savedValueItems = useMemo((): PermissionItem[] => {
		console.log('permissionRequest', permissionRequest);

		if(!permissionRequest.data?.roles?.length) {
			return [];
		}
		return permissionRequest.data!.roles[0].permissions?.map((i: Permission): PermissionItem => ({ 
			entityName: i.entity,
			operation: i.operation,
			status: 'saved',
		})) ?? []
	}, [permissionRequest.data]);

	const permissionValue = useMemo(() => {
		return groupByMap(savedValueItems,
			item => item.entityName,
			(item): Operation => ({ 
				operation: item.operation,
				status: item.status
			})
		);
	}, [savedValueItems])

	const id = useMemo(() => {
		return props.defaultValue?.id;
	}, [props.defaultValue]);

	const displayValue = useMemo(() => {
		if (!fullEntity) {
			return null;
		}

		const displayFieldProperty = fullEntity?.displayField?.name;
		if (!displayFieldProperty) {
			return null;
		}

		if (displayFieldProperty == 'id') {
			return null;
		}

		return props.defaultValue?.[displayFieldProperty];
	}, [fullEntity, props.defaultValue]);

	const upperValueAsFieldValues = useMemo(() => {
		return Object.fromEntries(
			Object.entries(myDialogContext.upperResults ?? {}).map(([key, value]) => {
				return [key, value];
			}))
	}, [myDialogContext.upperResults]);

	const formMethods = useForm({
		resolver: zodResolver(formSchema.schema),
		mode: 'all',
		defaultValues: {}
	});

	const onSave = useCallback((formData: any) => {
		myDialogContext.closeWithResults(formData);
	}, [myDialogContext.closeWithResults]);

	useImperativeHandle(ref, () => ({
		enterPressed: () => {
			formMethods.handleSubmit(onSave)();
		}
	}));

	const onPermissionValuesChanged = useCallback((entityName: string, operations: string[]) => {
		console.log(entityName, operations);
	}, [myDialogContext.closeWithResults]);

	return (
		<>
			<DynamicDialogHeader>
				Role
			</DynamicDialogHeader>
			<DialogContent>
				<Stack
					gap={2}>

					{!!id && <Id value={id} rootProps={{
						marginLeft: 'auto',
					}} />}

					<Form methods={formMethods} onSubmit={formMethods.handleSubmit(console.log)} >
						<Stack gap={1.5}>
							<TextShortFormField name="name" label="Name" required /> 
						</Stack>

						<PermissionSection values={permissionValue} valuesChanged={onPermissionValuesChanged} />
					</Form>

					{myDialogContext.openMode !== FormOpenMode.New && <BasicAuditTrail defaultValues={props.defaultValue} />}

				</Stack>
			</DialogContent>
			<DynamicDialogFooter>
				<Stack
					direction="row"
					gap={1}>
					<Button
						color="primary"
						variant="contained"
						onClick={formMethods.handleSubmit(onSave)}
					>Save</Button>
				</Stack>
			</DynamicDialogFooter>
		</>
	)
});

type PermissionItem = {
	entityName: string;
	operation: string;
	status: EdgeStatus;
}