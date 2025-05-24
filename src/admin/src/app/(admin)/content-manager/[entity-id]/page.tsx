"use client";

import PlusIcon from "@heroicons/react/24/outline/PlusIcon";
import { Button, ButtonGroup, Container, IconButton, Link, ListItemIcon, ListItemText, MenuItem, Stack, SvgIcon, Typography } from "@mui/material";
import { useCallback, useMemo, useRef, useState } from "react";
import { useDocumentTitle } from "@/hooks/use-document-title";
import { useEntities } from "@/hooks/use-entities";
import Avvvatars from 'avvvatars-react'
import { useDynamicDialog } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog";
import { DynamicDialog } from "@/core-features/dynamic-dialog/src/dynamic-dialog";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ContentManagerProvider, useContentManagerContext } from "@/features/content-manager/use-content-manager-context";
import { ContentManagerEntryDialogContent } from "@/features/content-manager/content-manager-entry-form";
import { ContentManagerScene } from "@/features/content-manager/content-manager-scene";
import { PageHeader } from "@/shared/components/page-header";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { SkeletonEntityPage } from "@/shared/components/skeleton-entity-page";
import { ActionsButton } from "@/shared/components/actions-button";
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { openOneFile, openOneJsonFile, saveJsonAs } from "@/lib/utils/open-file";
import { addNotification } from "@/lib/notifications/store";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { MenuButton } from "@/shared/components/menu/menu-button";

export default function EntityTypeDesigner() {

    useDocumentTitle({ title: 'Entity Type Designer' });
    const entityId = useCurrentEntityNameContext();
    const appStatus = useAppStatus();

    if(appStatus.isEntityFullyAvailable(entityId!) == false){
        return (<>
            <Container
                maxWidth={false}
                sx={{
                    height: '100%',
                    py: 2
                }}>

                <Stack spacing={2}>
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

                <Stack
                    spacing={2}>

                    <ContentManagerProvider entityName={entityId!}>
                        <PageContent entityId={entityId!} />
                    </ContentManagerProvider>

                </Stack>
            </Container>
        </>
    )
}

interface PageContentProps {
    entityId: string;
}
const PageContent = (props: PageContentProps) => {

    const graphEntities = useEntities();
    const entity = useMemo(() => graphEntities?.entities.find(i => i.name == props.entityId), [graphEntities?.entities, props.entityId]);

    const contentManager = useContentManagerContext();
    const dynamicDialog = useDynamicDialog();

    const addNewEntry = useCallback(() => {
        dynamicDialog.addPopup(ContentManagerEntryDialogContent, { entityName: contentManager.entityName }, FormOpenMode.New);
    }, [dynamicDialog, contentManager.entityName]);

    const finishDialog = useCallback(({ data, openMode, editId }: { data: any, openMode?: FormOpenMode, editId?: string }) => {

        if (openMode == FormOpenMode.Edit) {
            return;
        }

        if (openMode == FormOpenMode.New) {
            contentManager.contentManagerStore.addItem(data);
            return;
        }

    }, [contentManager.contentManagerStore]);

    const importData = useCallback(() => {
        openOneJsonFile().then((data: any[]) => {
            contentManager.contentManagerStore.addItems(data);
        }).catch(_ => {
            addNotification({ message: 'Invalid file', type: 'error' });
        });
    }, [contentManager.contentManagerStore.addItem]);

    const exportData = useCallback(() => {

        if(!contentManager.contentManagerStore?.state?.data || !contentManager.contentManagerStore?.state?.data?.length){
            addNotification({ message: 'No data to export', type: 'warning' });
            return;
        }
        saveJsonAs(`${ entity?.caption }.json`, contentManager.contentManagerStore.state.data);
    }, [contentManager.contentManagerStore?.state?.data]);

    if (!entity) {
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

                        <Avvvatars value={entity.caption!} style="shape" size={35} />
                        <Typography variant="h1">
                            {entity.caption}
                        </Typography>

                    </Stack>

                    <Stack
                        direction="row"
                        spacing={1}>

                        <Button
                            size="medium"
                            variant="contained"
                            startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                            onClick={() => addNewEntry()}
                            aria-label="Add new entry"
                            aria-haspopup="dialog">
                            Add New
                        </Button>


                        <MenuButton
                            slots={{
                                button: (<IconButton><MoreVertIcon /></IconButton>)
                            }}
                        >
                            <MenuItem data-autoclose onClick={() => importData()}>
                                <ListItemIcon>
                                    <CloudUploadIcon fontSize="small" />
                                </ListItemIcon>
                                <ListItemText>Import JSON</ListItemText>
                            </MenuItem>

                            <MenuItem data-autoclose onClick={() => exportData()}>
                                <ListItemIcon>
                                    <CloudDownloadIcon fontSize="small" />
                                </ListItemIcon>
                                <ListItemText>Export current page</ListItemText>
                            </MenuItem>

                        </MenuButton>
                    </Stack>

                </Stack>
            </PageHeader>

            <ContentManagerScene />

            <DynamicDialog
                ref={dynamicDialog.ref}
                finish={finishDialog} />
        </>
    );
}