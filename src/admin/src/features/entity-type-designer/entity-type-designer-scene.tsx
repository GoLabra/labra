import { Card, CardContent, Divider, ListItemIcon, MenuItem, SvgIcon, Tooltip } from "@mui/material";
import { useCallback, useMemo } from "react";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import PlusCircleIcon from "@heroicons/react/24/outline/PlusCircleIcon";
import EntityFieldType from "./entity-field-type-descriptor";
import { DynamicDialog } from "../../core-features/dynamic-dialog/src/dynamic-dialog";
import { useEntityTypeDesignerEntryDialog } from "./entity-type-designer-new-field-type-dialog";
import { ChildTypeDescriptor, getChildTypeForEntityChild, getIconForEntityChild } from "./designer-field-map";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { useDialog } from "@/hooks/use-dialog";
import { EntityTypeDesignerFieldFlags } from "./entity-type-designer-field-flags";
import StarBorderIcon from '@mui/icons-material/StarBorder';
import StarIcon from '@mui/icons-material/Star';
import { useEntityDesignerForEntity } from "./use-designer-entities";
import { Action, ColumnsFillRowSpacePlugin, CustomBodyCellContentRenderPlugin, EmptyDataPlugin, ColumnDef, HighlightRowPlugin, MosaicDataTable, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, SkeletonLoadingPlugin, useGridPlugins, usePluginWithParams } from "mosaic-data-table";
import { EmptyMessage } from "@/shared/components/empty-message";
import { ActionList } from "@/shared/components/action-list";
import { ActionListItem } from "@/shared/components/action-list-item";
import { ConfirmationDialog } from "@/shared/components/confirmation-dialog";
import { MuiCardFooter } from "@/shared/components/mui-card-footer";
import { DesignerEdge, DesignerField } from "@/types/entity";
import Defaults from "@/config/Defaults.json";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { SYSTEM_CHILDREN } from "@/config/CONST";
import { EntityOwner } from "@/lib/apollo/graphql.entities";

interface EntityTypeDesignerSceneProps {
    entityName: string;
}
export default function EntityTypeDesignerScene(
    props: EntityTypeDesignerSceneProps,
) {
    const { entityName } = props;

    const entityDesigner = useEntityDesignerForEntity(entityName);
    const deleteChildConfirmationDialog = useDialog<DesignerField | DesignerEdge>();
    const entityTypeDesignerEntryDialog = useEntityTypeDesignerEntryDialog();
    const appStatus = useAppStatus();
    const isSystemEntity = useMemo(() => entityDesigner.originalFullEntity?.owner == EntityOwner.Admin, [entityDesigner.originalFullEntity?.owner]);
    const isEntityBusy = useMemo(() => appStatus.isEntityInCodeGeneration(entityName), [appStatus.isEntityInCodeGeneration, entityName]);

    const onEntityChildrenDynamicDialogFinish = useCallback(({ data, openMode, editId }: {
        data: {
            selectedChildTypeDescriptor: ChildTypeDescriptor;
            field: DesignerField | DesignerEdge;
        };
        openMode?: FormOpenMode;
        editId?: string;
    }) => {
        if (openMode == FormOpenMode.New) {
            entityDesigner.addChild(data.field);
            return;
        }

        if (openMode == FormOpenMode.Edit) {
            if (!editId) {
                return;
            }
            entityDesigner.updateChild(editId, data.field);
            return;
        }
    }, [entityDesigner.addChild, entityDesigner.updateChild, entityName]);


    // Head Cells
    const columns = useMemo((): ColumnDef<DesignerField | DesignerEdge>[] => [
        {
            id: "caption",
            header: "Caption",
            cell: (child: DesignerField | DesignerEdge) => (
                <EntityFieldType
                    icon={getIconForEntityChild(child)}
                    label={child.caption}
                    designerStatus={child.designerStatus}
                    end={<EntityTypeDesignerFieldFlags child={child} />}
                ></EntityFieldType>
            ),

        },
        {
            id: "name",
            header: "Type",
            cell: (child: DesignerField | DesignerEdge) =>
                getChildTypeForEntityChild(child),

        },
        {
            id: "defaultValue",
            header: "",
            cell: (child: DesignerField | DesignerEdge) => (<>
                {entityDesigner.fullDesignerEntity.displayFieldCaption == child.caption && (
                    <Tooltip title="Used as Display Field" placement='bottom' arrow>
                        <SvgIcon fontSize="small" sx={{ verticalAlign: 'middle' }}>
                            <StarIcon fontSize="inherit" />
                        </SvgIcon>
                    </Tooltip>
                )}
            </>)
        }
    ], [entityDesigner.fullDesignerEntity?.displayFieldCaption]);

    // Row Actions
    const actions: Action<DesignerField | DesignerEdge>[] = useMemo(
        () => [
            {
                id: "edit",
                isVisible: (field: DesignerField | DesignerEdge) => SYSTEM_CHILDREN.includes(field.name) == false,
                render: (field: DesignerField | DesignerEdge) => (
                    <MenuItem
                        id="edit-menu-item"
                        key={`edit-${field.name}`}
                        disabled={isEntityBusy}
                        onClick={(e) => entityTypeDesignerEntryDialog.openEditChild(field)}
                    >
                        <ListItemIcon>
                            <EditIcon />
                        </ListItemIcon>
                        Edit
                    </MenuItem>
                ),
            }, {
                id: "displayField",
                render: (field: DesignerField | DesignerEdge) => (
                    <MenuItem
                        id="remove-menu-item"
                        key={`default-${field.name}`}
                        disabled={isEntityBusy}
                        onClick={(e) => entityDesigner.setDisplayField(field.caption)}
                    >
                        <ListItemIcon>
                            <StarBorderIcon />
                        </ListItemIcon>
                        Set as Display Field
                    </MenuItem>
                ),
                isVisible: (child: DesignerField | DesignerEdge) => {
                    if (child.__typename === "Edge") {
                        return false;
                    }
                    return child.caption != entityDesigner.fullDesignerEntity.displayFieldCaption;
                }
            }, {
                id: "r-divider",
                isVisible: (field: DesignerField | DesignerEdge) => SYSTEM_CHILDREN.includes(field.name) == false,
                render: (field: DesignerField | DesignerEdge) => <Divider key={`divider1-${field.name}`} />
            }, {
                id: "remove",
                isVisible: (field: DesignerField | DesignerEdge) => SYSTEM_CHILDREN.includes(field.name) == false,
                render: (field: DesignerField | DesignerEdge) => (
                    <MenuItem
                        id="remove-menu-item"
                        key={`remove-${field.caption}`}
                        disabled={isEntityBusy}
                        onClick={(e) =>
                            deleteChildConfirmationDialog.handleOpen(field)
                        }
                    >
                        <ListItemIcon>
                            <DeleteIcon />
                        </ListItemIcon>
                        Remove
                    </MenuItem>
                ),
            }
        ], [entityTypeDesignerEntryDialog.openEditChild, deleteChildConfirmationDialog, entityName, isEntityBusy, entityTypeDesignerEntryDialog, entityDesigner.setDisplayField, entityDesigner.fullDesignerEntity?.displayFieldCaption]);

    const gridPlugins = useGridPlugins(
        CustomBodyCellContentRenderPlugin,
        usePluginWithParams(PaddingPluggin, {}),
        ColumnsFillRowSpacePlugin,
        usePluginWithParams(RowActionsPlugin, {
            actions: actions,
            visible: isSystemEntity == false
        }),
        usePluginWithParams(HighlightRowPlugin, {
            isRowHighlighted: useCallback((row: any) => row.caption === entityDesigner.fullDesignerEntity.displayFieldCaption, [entityDesigner.fullDesignerEntity?.displayFieldCaption]),
        }),
        PinnedColumnsPlugin,
        usePluginWithParams(SkeletonLoadingPlugin, {
            isLoading: entityDesigner.fullDesignerEntity?.loading || isEntityBusy,
            rowsWhenEmpty: Defaults.dataTable.skeletonRowsCount
        }),
        usePluginWithParams(EmptyDataPlugin, {
            content: <EmptyMessage />
        }),
    );

    return (<>
        <Card>
            <CardContent>
                <MosaicDataTable
                    plugins={gridPlugins}
                    caption="Entity fields and edges"
                    items={entityDesigner.fullDesignerEntity?.children.filter(i => i.designerStatus != 'deleted')}
                    headCells={columns}
                />
            </CardContent>


            {isSystemEntity == false && (<>
                <Divider />

                <MuiCardFooter>
                    <ActionList>
                        <ActionListItem
                            onClick={() => entityTypeDesignerEntryDialog.openAddNewChild()}
                            disabled={isEntityBusy}
                            icon={
                                <SvgIcon fontSize="small">
                                    <PlusCircleIcon />
                                </SvgIcon>
                            }
                            label="Add New"
                        />
                    </ActionList>
                </MuiCardFooter>
            </>)}
        </Card>

        <DynamicDialog
            ref={entityTypeDesignerEntryDialog.ref}
            finish={onEntityChildrenDynamicDialogFinish} />

        <ConfirmationDialog
            message="Are you sure you want to delete this field?"
            onCancel={deleteChildConfirmationDialog.handleClose}
            onConfirm={() => {
                deleteChildConfirmationDialog.handleClose();
                entityDesigner.deleteChild(deleteChildConfirmationDialog.data!.name);
            }}
            open={deleteChildConfirmationDialog.open}
            title="Delete Confirmation"
            variant="error"
        />
    </>);
}
EntityTypeDesignerScene.displayName = 'EntityTypeDesignerScene';