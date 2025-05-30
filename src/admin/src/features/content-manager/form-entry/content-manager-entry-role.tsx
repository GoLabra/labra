import { ChainDialogContentRef } from "@/core-features/dynamic-dialog/src/dynamic-dialog-types";
import { forwardRef, Fragment, use, useCallback, useEffect, useImperativeHandle, useMemo } from "react";
import { ContentManagerEntryDialogContentProps } from "../content-manager-entry-form";
import { DynamicDialogFooter } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-footer";
import { DynamicDialogHeader } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog-header";
import { useMyDialogContext } from "@/core-features/dynamic-dialog/src/use-my-dialog-context";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { BasicAuditTrail } from "@/shared/components/basic-audit-trail";
import { Id } from "@/shared/components/id";
import { zodResolver } from "@hookform/resolvers/zod";
import { DialogContent, Stack, Button, Typography, Checkbox, Box, Divider, Chip, Tooltip, MenuItem, ListItemIcon, ListItemText, IconButton } from "@mui/material";
import { get, useForm, useFormContext } from "react-hook-form";
import { useContentManagerFormSchema, useGenericEdgeFormSchema, useGenericFiedFormSchema } from "../use-content-manager-form-schema";
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
import React from "react";
import { MenuButton } from "@/shared/components/menu/menu-button";
import MoreVertIcon from '@mui/icons-material/MoreVert';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';

const GET_ROLE_PERMISSION_QUERY = gql`query getRolePermissionQuery($where: RoleWhereInput) {
	roles(where: $where)  {
		name
		permissions {
			id
    		entity
    		operation
    	}
  	}
}`


type Operation = {
	id?: string;
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

const changePermission = (allPermissions: Record<string, Operation[]>, entityName: string, operations: string[]) => {

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

	return {
		[entityName]: [
			...save,
			...create,
			...disconnect
		]
	}
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

	const isAllChecked = useMemo(() => props.operations.length == props.operationDefs.length, [props.operations, props.operationDefs]);

	const allChanged = useCallback((event: GenericEvent<boolean>) => {
		if (event.target.value) {
			props.valuesChanged?.(props.operationDefs.map(i => i.name));
			return;
		}

		props.valuesChanged?.([]);
	}, [props.valuesChanged]);

	return (
		<Stack gap={1.5} my={2} borderRadius={0}>
			<Typography >
				{props.entityCaption}
			</Typography>

			<Stack direction="row" gap={1} flexWrap="wrap" justifyContent="space-evenly" width={1}>

				<BooleanFormComponent
						name={`${props.entityCaption}-all`}
						label="All"
						color="success"
						value={isAllChecked}
						onChange={allChanged}
					/>

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
		const newPermissions = changePermission(allPermissions, entityName, operations);

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: {
					...allPermissions,
					...newPermissions
				}
			}
		});
	}, [props.name, formControllerHandler.value, formControllerHandler.onChange]);

	const selectAll = useCallback(() => {

		const allPermissions = formControllerHandler.value ?? {};

		const newPermissions = entities.reduce((acc: any, entity) => {
			return {
				...acc,
				...changePermission(allPermissions, entity.name, operationDefs.map(i => i.name))
			}
		}, {});

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: {
					...allPermissions,
					...newPermissions
				}
			}
		});
	}, [entities, onPermissionValuesChanged]);

	const permissions = useMemo(() => {
		return Object.values(formControllerHandler.value).reduce((sum, arr) => [...sum, ...arr], []);
	}, [formControllerHandler.value]);

	const permissionsCount = useMemo(() => {
		return permissions.filter(i => i.status != 'disconnect');
	}, [permissions]);

	const tooltipLabel = useMemo(() => {
		const added = permissions.filter(i => i.status == 'create').length;
		const disconnect = permissions.filter(i => i.status == 'disconnect').length;
		const saved = permissions.filter(i => i.status == 'saved').length;

		return `${saved} existing, ${added} added, ${disconnect} removed`;
	}, [permissions]);

	return (
		<Stack gap={1.5}>

			<SectionTitle title="Permissions"
				endAdornment={
					<MenuButton
						slots={{
							button: (<IconButton><MoreVertIcon /></IconButton>)
						}}
					>
						<MenuItem data-autoclose onClick={selectAll}>
							<ListItemIcon>
								<CheckCircleIcon fontSize="small" />
							</ListItemIcon>
							<ListItemText>Select All Permissions</ListItemText>
						</MenuItem>


					</MenuButton>
				}>
				<Tooltip title={tooltipLabel} placement='bottom' arrow>
					<Chip component="span" label={permissionsCount.length} color="secondary" variant="outlined" size="small" />
				</Tooltip>


			</SectionTitle>

			{entities.map((entity) => {
				return <Fragment key={entity.name}>
					<EntityPermissionView
						entityCaption={entity.caption}
						operationDefs={operationDefs}
						operations={(formControllerHandler.value?.[entity.name] ?? []).filter(i => i.status !== 'disconnect')
							.map(i => i.operation) ?? []}
						valuesChanged={(values) => {
							onPermissionValuesChanged(entity.name, values);
						}} />

					<Divider />
				</Fragment>
			})}

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
			disconnect: allEntityPermissions.filter(i => i.status == 'disconnect')
				.filter(i => i.id)
				.map(i => ({ id: i.id }))
		}
	})
});

export const ContentManagerEntryRole = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	const myDialogContext = useMyDialogContext();
	const fullEntity = useFullEntity({ entityName: props.entityName });
	const permissionRequest = useQuery<{ roles: Role[] }>(GET_ROLE_PERMISSION_QUERY, {
		variables: {
			where: {
				id: myDialogContext.editId
			}
		},
		fetchPolicy: 'network-only',
		skip: myDialogContext.openMode === FormOpenMode.New
	});

	const permissionValue = useMemo(() => {

		if (!permissionRequest.data?.roles?.length) {
			return [];
		}

		const savedValueItems = permissionRequest.data!.roles[0].permissions?.map((i: Permission): PermissionItem => ({
			id: i.id,
			entityName: i.entity,
			operation: i.operation,
			status: 'saved',
		})) ?? [];

		if (!savedValueItems.length) {
			return null;
		}

		return groupByMap(savedValueItems,
			item => item.entityName,
			(item): Operation => ({
				id: item.id,
				operation: item.operation,
				status: item.status
			})
		);

	}, [permissionRequest.data]);


	const id = useMemo(() => {
		return props.defaultValue?.id;
	}, [props.defaultValue]);

	const usersRoles = useGenericEdgeFormSchema(fullEntity, 'userRoles');

	const computedSchema = useMemo(() => {
		return schema.merge(z.object({
			userRoles: usersRoles.schema
		}));
	}, [schema, usersRoles.schema]);

	const computedDefaultValues = useMemo(() => {
		return {
			name: props.defaultValue?.name,
			userRoles: usersRoles.convertFromRawValue(props.defaultValue?.userRoles),
			permissions: permissionValue
		}
	}, [props.defaultValue, permissionValue]);

	const formMethods = useForm({
		resolver: zodResolver(computedSchema),
		mode: 'all',
		defaultValues: computedDefaultValues
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
							{usersRoles.field}
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
	id?: string;
	entityName: string;
	operation: string;
	status: EdgeStatus;
}