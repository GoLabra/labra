"use client";

import PlusIcon from "@heroicons/react/24/outline/PlusIcon";
import { Container, ListItemIcon, ListItemText, MenuItem, Stack, SvgIcon, Typography } from "@mui/material";
import { useCallback } from "react";
import { useDocumentTitle } from "@/hooks/use-document-title";
import Avvvatars from 'avvvatars-react'
import { useDynamicDialog } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog";
import { DynamicDialog, FinishResult } from "@/core-features/dynamic-dialog/src/dynamic-dialog";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { ContentManagerProvider, useContentManagerContext } from "@/features/content-manager/use-content-manager-context";
import { ContentManagerEntryDialogContent } from "@/features/content-manager/content-manager-entry-form";
import { ContentManagerGrid } from "@/features/content-manager/content-manager.grid";
import { PageHeader } from "@/shared/components/page-header";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { SkeletonEntityPage } from "@/shared/components/skeleton-entity-page";
import CloudDownloadIcon from '@mui/icons-material/CloudDownload';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import { openOneJsonFile, saveJsonAs } from "@/lib/utils/open-file";
import { addNotification } from "@/lib/notifications/store";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import MoreVertIcon from '@mui/icons-material/MoreVert';
import { MenuButton } from "@/shared/components/menu/menu-button";
import { ShortcutButton, ShortcutIconButton } from "@/shared/components/key-handler/with-click-shortcut";
import { Key } from "@/shared/components/key-handler/types";

export default function ContentManager() {

	useDocumentTitle({ title: 'Entity Type Designer' });
	const entityId = useCurrentEntityNameContext();
	const isEntityFullyAvailable = useAppStatus().isEntityFullyAvailable(entityId!);

	return (
		<Container
			maxWidth={false}
			sx={{
				height: '100%',
				py: 2
			}}>

			<Stack spacing={2}>
				<ContentManagerProvider entityName={entityId!}>
					{isEntityFullyAvailable == false && <SkeletonEntityPage />}
					{isEntityFullyAvailable == true && <PageContent />}
				</ContentManagerProvider>
			</Stack>
		</Container>
	)
}

const PageContent = () => {

	const contentManager = useContentManagerContext();
	const dynamicDialog = useDynamicDialog();

	const addNewEntry = useCallback(() => {
		dynamicDialog.addPopup(ContentManagerEntryDialogContent, { entityName: contentManager.entityName }, FormOpenMode.New);
	}, [dynamicDialog, contentManager.entityName]);

	const finishDialog = useCallback(({ data, openMode }: FinishResult) => {

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

		if (!contentManager.contentManagerStore?.state?.data || !contentManager.contentManagerStore?.state?.data?.length) {
			addNotification({ message: 'No data to export', type: 'warning' });
			return;
		}
		saveJsonAs(`${contentManager.fullEntity?.caption}.json`, contentManager.contentManagerStore.state.data);
	}, [contentManager.contentManagerStore?.state?.data]);

	if (!contentManager.fullEntity) {
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
						<Avvvatars value={contentManager.fullEntity?.caption!} style="shape" size={35} />
						<Typography variant="h1">
							{contentManager.fullEntity?.caption}
						</Typography>
					</Stack>

					<Stack
						direction="row"
						alignItems="center"
						spacing={1}>

						<ShortcutButton
							shortcutKey={Key.n}
							iconOpacity={true}
							size="medium"
							variant="contained"
							startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
							onClick={() => addNewEntry()}
							aria-label="Add new entry"
							aria-haspopup="dialog">
							Add New
						</ShortcutButton>

						<MenuButton
							slots={{
								button: (<ShortcutIconButton shortcutKey={Key.o} tooltip="Options"><MoreVertIcon /></ShortcutIconButton>)
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

			<ContentManagerGrid />

			<DynamicDialog
				ref={dynamicDialog.ref}
				finish={finishDialog} />
		</>
	);
}