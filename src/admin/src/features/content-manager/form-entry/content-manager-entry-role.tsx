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
import HistoryIcon from "@mui/icons-material/History";
import RadioButtonUncheckedIcon from '@mui/icons-material/RadioButtonUnchecked';
import { useRelationDiff } from "../use-relation-diff";

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

type PermissionItem = {
	id?: string;
	entityName: string;
	operation: string;
	status: EdgeStatus;
}

type Operation = {
	id?: string;
	operation: string;
	status: EdgeStatus;
}

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


type OperationDef = {
	name: string;
	caption: string;
}

const changePermission = (savedPermissions: PermissionItem[], entityName: string, operations: string[]): PermissionItem[] => {

	const currentEntityPermissions = savedPermissions.filter(i => i.entityName == entityName);

	const disconnect = currentEntityPermissions.filter(i => operations.includes(i.operation) == false)
		.map((i): PermissionItem => ({
			...i,
			status: 'delete'
		}));
		
	const create = operations.filter(i => !!savedPermissions.find(j => j.entityName == entityName && j.operation == i) == false)
		.map((i): PermissionItem => ({
			entityName: entityName,
			operation: i,
			status: 'create'
		}));

	return [
		...create,
		...disconnect
	]
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
	const formControllerHandler = useLiteController<PermissionItem[]>({ name: props.name, control: formContext.control });
	const myDialogContext = useMyDialogContext();

	
	
	const permissionRequest = useQuery<{ roles: Role[] }>(GET_ROLE_PERMISSION_QUERY, {
		variables: {
			where: {
				id: myDialogContext.editId
			}
		},
		fetchPolicy: 'network-only',
		skip: myDialogContext.openMode === FormOpenMode.New
	});


	const saved = useMemo(() => {

		if (!permissionRequest.data?.roles?.length) {
			return [];
		}

		const savedValueItems = permissionRequest.data!.roles[0].permissions?.map((i: Permission): PermissionItem => ({
			id: i.id,
			entityName: i.entity,
			operation: i.operation,
			status: 'saved',
		}));

		return savedValueItems ?? [];
	}, [permissionRequest.data]);


	const relationDiff = useRelationDiff<PermissionItem>({ saved, changedArray: formControllerHandler.value });

	// const currentPermissions = useMemo((): PermissionItem[] => {
	// 	const all = [...existing, ...formControllerHandler.value ?? []];
	// 	return all.filter(i => {
	// 		if (i.status == 'delete') {
	// 			return false;
	// 		}

	// 		if (i.status == 'saved') {
	// 			if (all.find(j => j.id == i.id && j.status == 'delete')) {
	// 				return false;
	// 			}
	// 		}

	// 		return true;
	// 	});
	// }, [existing, formControllerHandler.value]);

	const allPermissionsGrouped = useMemo((): Record<string, Operation[]> => {

		const permissionByEntity = groupByMap(relationDiff.showingItems,
			item => item.entityName,
			(item): Operation => ({
				id: item.id,
				operation: item.operation,
				status: item.status
			})
		);

		return permissionByEntity;

	}, [relationDiff.showingItems]);


	const onPermissionValuesChanged = useCallback((entityName: string, operations: string[]) => {
		const newPermissions = changePermission(saved, entityName, operations);

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: [
					...formControllerHandler.value?.filter(i => i.entityName != entityName) ?? [],
					...newPermissions
				]
			}
		});
	}, [saved, formControllerHandler.value, formControllerHandler.onChange]);

	const selectAll = useCallback(() => {
		const newPermissions = entities.flatMap((entity) => {
			return changePermission(saved, entity.name, operationDefs.map(i => i.name))
		});

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: newPermissions
			}
		});
	}, [entities, onPermissionValuesChanged]);

	const clearAll = useCallback(() => {
		const newPermissions = entities.flatMap((entity) => {
			return changePermission(saved, entity.name, [])
		});

		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: newPermissions
			}
		});
	}, [entities, onPermissionValuesChanged]);

	const resetChanges = useCallback(() => {
		formControllerHandler.onChange({
			target: {
				name: props.name,
				value: []
			}
		});
	}, [entities, onPermissionValuesChanged]);

	const permissionsCount = useMemo(() => {
		return relationDiff.showingItems.filter(i => i.status != 'delete');
	}, [relationDiff.showingItems]);

	const tooltipLabel = useMemo(() => {
		const itemsAdded = formControllerHandler.value?.filter(i => i.status == 'create').length ?? 0;
		const itemsDisconnected = formControllerHandler.value?.filter(i => i.status == 'delete').length ?? 0;
		const itemsSaved = saved.filter(i => formControllerHandler.value?.find(j => j.id == i.id && j.status == 'delete') == null).length ?? 0;

		return `${itemsSaved} existing, ${itemsAdded} added, ${itemsDisconnected} removed`;
	}, [saved, formControllerHandler.value]);

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
						<MenuItem data-autoclose onClick={clearAll}>
							<ListItemIcon>
								<RadioButtonUncheckedIcon fontSize="small" />
							</ListItemIcon>
							<ListItemText>Clear All Permissions</ListItemText>
						</MenuItem>

						{myDialogContext.openMode === FormOpenMode.Edit && (<MenuItem data-autoclose onClick={resetChanges}>
							<ListItemIcon>
								<HistoryIcon fontSize="small" />
							</ListItemIcon>
							<ListItemText>Reset Changes</ListItemText>
						</MenuItem>)}

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
						operations={(allPermissionsGrouped[entity.name] ?? []).filter(i => i.status !== 'delete')
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
	permissions: z.any().optional().nullish().transform((val: PermissionItem[] | undefined) => {
		if (!val) {
			return val;
		}

		return {
			connect: val.filter(i => i.status == 'connect')
				.map(i => ({
					entity: i.entityName,
					operation: i.operation
				})),
			create: val.filter(i => i.status == 'create')
				.map(i => ({
					entity: i.entityName,
					operation: i.operation
				})),
			delete: val.filter(i => i.status == 'delete')
				.filter(i => i.id)
				.map(i => ({ id: i.id }))
		}
	})
});

export const ContentManagerEntryRole = forwardRef<ChainDialogContentRef, ContentManagerEntryDialogContentProps>((props, ref) => {

	const myDialogContext = useMyDialogContext();
	const fullEntity = useFullEntity({ entityName: props.entityName });

	const id = useMemo(() => {

		if (myDialogContext.openMode === FormOpenMode.New) {
			return null
		}
		return myDialogContext.editId;
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
		}
	}, [props.defaultValue]);

	const formMethods = useForm({
		resolver: zodResolver(computedSchema),
		mode: 'all',
		defaultValues: computedDefaultValues
	});

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

ContentManagerEntryRole.displayName = 'ContentManagerEntryRole';
