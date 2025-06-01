"use client";

import { Button, ButtonGroup, Container, Divider, IconButton, ListItemIcon, ListItemText, MenuItem, Stack, SvgIcon, Typography } from "@mui/material";
import { useCallback, useMemo } from "react";
import { useParams } from "next/navigation";
import { DynamicDialog, FinishResult } from "@/core-features/dynamic-dialog/src/dynamic-dialog";
import { useDialog } from "@/hooks/use-dialog";
import { VscSave } from "react-icons/vsc";
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import AddIcon from '@mui/icons-material/Add';
import HistoryIcon from '@mui/icons-material/History';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { useDocumentTitle } from "@/hooks/use-document-title";
import { useDynamicDialog } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog";
import { AppStatus, Entity, EntityOwner } from "@/lib/apollo/graphql.entities";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import Avvvatars from 'avvvatars-react'
import { EntityTypeChanges } from "@/features/entity-type-designer/entity-type-changes";
import { EntityTypeDesignerEntryDialogContent } from "@/features/entity-type-designer/entity-type-designer-new-field-type-dialog";
import EntityTypeDesignerScene from "@/features/entity-type-designer/entity-type-designer-scene";
import { EntityTypeNewEntityDialog } from "@/features/entity-type-designer/entity-type-new-entity-dialog";
import { useEntitiesDesigner, useEntityDesignerForEntity } from "@/features/entity-type-designer/use-designer-entities";
import { ActionsButton } from "@/shared/components/actions-button";
import { PageHeader } from "@/shared/components/page-header";
import { ConfirmationDialog } from "@/shared/components/confirmation-dialog";
import { ChangedNameCaptionEntity, DesignerEdge, DesignerField } from "@/types/entity";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { SkeletonEntityPage } from "@/shared/components/skeleton-entity-page";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import { MenuButton } from "@/shared/components/menu/menu-button";
import { ChildTypeDescriptor } from "@/features/entity-type-designer/designer-field-map";

export default function EntityTypeDesigner() {

    useDocumentTitle({ title: 'Entity Type Designer' });
    const entityId = useCurrentEntityNameContext();
    const appStatus = useAppStatus();

    if (appStatus.isEntityFullyAvailable(entityId!) == false) {
        return (<>
            <Container
                maxWidth={false}
                sx={{
                    height: '100%',
                    py: 2
                }}>

                <Stack
                    spacing={2}>

                    <SkeletonEntityPage />

                </Stack>
            </Container>
        </>)
    }

    return (
        <>
            <Container
                maxWidth={false}
                sx={{
                    height: '100%',
                    py: 2
                }}>

                <Stack spacing={2}>
                    <PageContent entityId={entityId!} />
                </Stack>
            </Container>

        </>
    )
}


interface PageContentProps {
    entityId: string;
}
const PageContent = (props: PageContentProps) => {

    const entityDesigner = useEntityDesignerForEntity(props.entityId);
    const isSystemEntity = useMemo(() => entityDesigner.originalFullEntity?.owner == EntityOwner.Admin, [entityDesigner.originalFullEntity?.owner]);
    const entityDynamicDialog = useDynamicDialog();
    const entityChildrenDynamicDialog = useDynamicDialog();
    const appStatus = useAppStatus();
    const saveChangesDialogConfirmation = useDialog();
    const deleteEntityConfirmationDialog = useDialog();
    const revertAllChangesConfirmationDialog = useDialog();

    const onEntityDynamicDialogFinish = useCallback(({ data, openMode, editId }: FinishResult<ChangedNameCaptionEntity>) => {
        if (!editId) {
            return;
        }

        entityDesigner.editEntity(data);
    }, [entityDesigner]);

    const onEntityChildrenDynamicDialogFinish = useCallback(({ data, openMode, editId }: {
        data: {
            selectedChildTypeDescriptor: ChildTypeDescriptor;
            field: DesignerField | DesignerEdge;
        };
        openMode?: FormOpenMode;
        editId?: string;
    }) => {
        if (openMode == FormOpenMode.Edit) {
            return;
        }
        // if (!entityDesigner.fullDesignerEntity) {
        //     return;
        // }
        entityDesigner.addChild(data.field);

    }, [entityDesigner.addChild]);

    if (!entityDesigner.fullDesignerEntity) {
        return null;
    }

    return (
        <>


            <PageHeader
                sx={{
                    pl: 1,
                    pr: 2
                }}>
                <Stack
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                    spacing={1}>

                    <Stack direction="row" alignItems="center" gap={1}>
                        <Avvvatars value={entityDesigner.fullDesignerEntity?.caption!} style="shape" size={35} />
                        <Typography variant="h2">
                            {entityDesigner.fullDesignerEntity?.caption!}
                        </Typography>
                    </Stack>

                    {isSystemEntity == false && (
                        <Stack
                            direction="row"
                            spacing={1}>

                            <Button
                                size="medium"
                                variant="contained"
                                startIcon={<VscSave size={16} />}
                                onClick={() => saveChangesDialogConfirmation.handleOpen()}
                                disabled={entityDesigner.hasChanges == false || appStatus.applicationStatus != AppStatus.Up}
                                aria-label="Save Changes"
                                aria-haspopup="dialog"
                            >
                                Save
                            </Button>


                            <MenuButton
                                slots={{
                                    button: (<IconButton><MoreVertIcon /></IconButton>)
                                }}
                            >
                                <MenuItem data-autoclose
                                    onClick={() => entityDynamicDialog.addPopup(EntityTypeNewEntityDialog, { defaultValue: entityDesigner.fullDesignerEntity }, FormOpenMode.Edit, entityDesigner.fullDesignerEntity.name)}>
                                    <ListItemIcon>
                                        <EditIcon fontSize="small" />
                                    </ListItemIcon>
                                    <ListItemText>Edit Entity</ListItemText>
                                </MenuItem>

                                <MenuItem data-autoclose
                                    onClick={() => entityChildrenDynamicDialog.addPopup(EntityTypeDesignerEntryDialogContent, {}, FormOpenMode.New)}>
                                    <ListItemIcon>
                                        <AddIcon fontSize="small" />
                                    </ListItemIcon>
                                    <ListItemText>Add New Field / Relation</ListItemText>
                                </MenuItem>


                                <Divider />

                                <MenuItem
                                    data-autoclose
                                    disabled={entityDesigner.hasChanges == false}
                                    onClick={() => revertAllChangesConfirmationDialog.handleOpen()}>
                                    <ListItemIcon>
                                        <HistoryIcon fontSize="small" />
                                    </ListItemIcon>
                                    <ListItemText>Revert all changes</ListItemText>
                                </MenuItem>

                                <Divider />

                                <MenuItem data-autoclose
                                    onClick={() => deleteEntityConfirmationDialog.handleOpen()}
                                    sx={{
                                        color: "error.main"
                                    }}>
                                    <ListItemIcon>
                                        <DeleteIcon fontSize="small" color="error" />
                                    </ListItemIcon>
                                    <ListItemText>Remove Entity</ListItemText>
                                </MenuItem>

                            </MenuButton>

                        </Stack>)}

                </Stack>
            </PageHeader>

            {entityDesigner.entityStatus != 'deleted' && (
                <EntityTypeDesignerScene entityName={props.entityId}></EntityTypeDesignerScene>
            )}

            <DynamicDialog
                ref={entityDynamicDialog.ref}
                finish={onEntityDynamicDialogFinish} />

            <DynamicDialog
                ref={entityChildrenDynamicDialog.ref}
                finish={onEntityChildrenDynamicDialogFinish} />

            <ConfirmationDialog
                message={<><EntityTypeChanges entityName={props.entityId} onClose={saveChangesDialogConfirmation.handleClose} /></>}
                onCancel={saveChangesDialogConfirmation.handleClose}
                onConfirm={() => { saveChangesDialogConfirmation.handleClose(); entityDesigner.saveChanges(); }}
                open={saveChangesDialogConfirmation.open}
                title="Save Entities Changes"
                variant="info"
                PaperProps={{
                    sx: {
                        width: '800px', // Set the desired width here
                        maxWidth: '800px' // Prevents the default maxWidth behavior
                    }
                }} />

            <ConfirmationDialog
                message="Are you sure you want to delete this entity?"
                onCancel={deleteEntityConfirmationDialog.handleClose}
                onConfirm={() => { deleteEntityConfirmationDialog.handleClose(); entityDesigner.deleteEntity(); }}
                open={deleteEntityConfirmationDialog.open}
                title="Delete Confirmation"
                variant="error" />

            <ConfirmationDialog
                message="Are you sure you want to revert all changes?"
                onCancel={revertAllChangesConfirmationDialog.handleClose}
                onConfirm={() => { revertAllChangesConfirmationDialog.handleClose(); entityDesigner.revertEntityChange(); }}
                open={revertAllChangesConfirmationDialog.open}
                title="Revert All Changes Confirmation"
                variant="error" />
        </>
    )
}