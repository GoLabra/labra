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
import { get, useForm, useFormContext } from "react-hook-form";
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
import { useLiteController } from "@/core-features/dynamic-form/lite-controller";
import { SavedSearch } from "@mui/icons-material";
import { permission } from "node:process";


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

// type EntityPermissions = {
// 	entityName: string;
// 	operation: Operation[];
// }

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
		if (event.target.value) {
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

interface PermissionSectionProps {
	name: string
}
const PermissionSection = (props: PermissionSectionProps) => {

	const { entities } = useEntities();
	const formContext = useFormContext();
	const formControllerHandler = useLiteController<Record<string, Operation[]>>({ name: props.name, control: formContext.control });

	const operationDefs = [{
		caption: 'Create',
		name: 'Create'
	}, {
		caption: 'Read',
		name: 'Read'
	}, {
		caption: 'Update',
		name: 'Update'
	}, {
		caption: 'Delete',
		name: 'Delete'
	}];

	const onPermissionValuesChanged = useCallback((entityName: string, operations: string[]) => {

		const allPermissions = formControllerHandler.value ?? {};
		const currentEntityPermissions = allPermissions[entityName] ?? [];

		const save = [
			...currentEntityPermissions.filter(i => operations.includes(i.operation))
				.filter(i => i.status === 'saved'),

			...currentEntityPermissions.filter(i => operations.includes(i.operation))
				.filter(i => i.status === 'disconnect')
				.map(i => ({
					...i,
					status: 'saved'
				}))
		]

		const disconnect = [
			...currentEntityPermissions.filter(i => i.status === 'saved')
				.filter(i => operations.includes(i.operation) == false)
				.map(i => ({
					...i,
					status: 'disconnect'
				})),
			...currentEntityPermissions.filter(i => i.status === 'disconnect')
				.filter(i => operations.includes(i.operation) == false),
		];

		const savedAndDisconnect = [
			...save.map(i => i.operation),
			...disconnect.map(i => i.operation)
		];

		const create = operations.filter(i => savedAndDisconnect.includes(i) == false)
			.map(i => ({
				operation: i,
				status: 'create'
			}));

		console.log([
			...save,
			...create,
			...disconnect
		]);

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: {
					...allPermissions,
					[entityName]: [
						...save,
						...create,
						...disconnect
					]
				}
			}
		});
	}, [props.name, formControllerHandler.value, formControllerHandler.onChange]);

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
						operations={(formControllerHandler.value?.[entity.name] ?? []).filter(i => i.status !== 'disconnect')
							.map(i => i.operation) ?? []}
						valuesChanged={(values) => {
							onPermissionValuesChanged(entity.name, values);
						}} />
				})}
			</Box>


		</Stack>
	)
}
// END OF PERMISSION SECTION



export const schema = z.object({
	name: z.string().min(1, 'Caption is required'),
	permissions: z.any().optional().nullish().transform((val: Record<string, Operation[]> | undefined) => {
		if (!val) {
			return val;
		}

		const allEntityPermissions = Object.keys(val)
											.flatMap((entityName) => {
												return val[entityName].map((operation) => ({
													entityName,
													...operation
												}))
											});
		
		return {
			connect: allEntityPermissions.filter(i => i.status == 'connect')
				.map(i => ({
					entity: i.entityName,
					operation: i.operation
				})),
			create: allEntityPermissions.filter(i => i.status == 'create')
				.map(i => ({
					entity: i.entityName,
					operation: i.operation
				})),
			// disconnect: allEntityPermissions.filter(i => i.status == 'disconnect')
			// 	.map(i => ({
			// 		entity: i.entityName,
			// 		operation: i.operation
			// 	}))
		}
	})
});

export const ContentManagerEntryRole = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	const myDialogContext = useMyDialogContext();

	const permissionRequest = useQuery<{ roles: Role[] }>(GET_ROLE_PERMISSION_QUERY, {
		variables: {
			where: {
				id: props.defaultValue?.id
			}
		},
		fetchPolicy: 'network-only',
		skip: myDialogContext.openMode === FormOpenMode.New
	});

	const savedValueItems = useMemo((): PermissionItem[] => {

		if (!permissionRequest.data?.roles?.length) {
			return [];
		}
		return permissionRequest.data!.roles[0].permissions?.map((i: Permission): PermissionItem => ({
			entityName: i.entity,
			operation: i.operation,
			status: 'saved',
		})) ?? []
	}, [permissionRequest.data]);

	const permissionValue = useMemo(() => {

		if (!savedValueItems.length) {
			return null;
		}

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

	const formMethods = useForm({
		resolver: zodResolver(schema),
		mode: 'all',
		defaultValues: {
			name: props.defaultValue?.name,
			permissions: permissionValue
		}
	});

	useEffect(() => {
		if (!permissionValue) {
			return;
		}

		formMethods.reset({
			...formMethods.getValues(), // Keep other form values
			permissions: permissionValue
		});
	}, [permissionValue]);

	const onSave = useCallback((formData: any) => {
		myDialogContext.closeWithResults(formData);
	}, [myDialogContext.closeWithResults]);

	useImperativeHandle(ref, () => ({
		enterPressed: () => {
			formMethods.handleSubmit(onSave)();
		}
	}));

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

						<PermissionSection name="permissions" />
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