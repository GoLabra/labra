import { Box, Card, CardContent, CardHeader, Divider, ListItemIcon, MenuItem, SvgIcon } from "@mui/material"
import PlusCircleIcon from "@heroicons/react/24/outline/PlusCircleIcon"
import { useCallback, useEffect, useMemo, useState } from "react"
import DeleteIcon from "@mui/icons-material/Delete";
import EditIcon from "@mui/icons-material/Edit";
import { ContentManagerSearch } from "./content-manager-search"
import { useSelection } from "@/hooks/use-selection";
import { DynamicDialog, FinishResult } from "../../core-features/dynamic-dialog/src/dynamic-dialog"
import { FormOpenMode } from "@/core-features/dynamic-form/form-field"
import { useDialog } from "@/hooks/use-dialog"
import { useContentManagerIds } from "./use-content-manger-ids"
import { useDynamicDialog } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog"
import { ContentManagerEntryDialogContent } from "./content-manager-entry-form"
import { Action, ColumnsFillRowSpacePlugin, ColumnSortPlugin, CustomBodyCellContentRenderPlugin, EmptyDataPlugin, ColumnDef, HighlightColumnPlugin, HighlightRowPlugin, MosaicDataTable, Order, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, RowExpansionPlugin, RowSelectionPlugin, SkeletonLoadingPlugin, useGridPlugins, usePluginWithParams, useRowExpansionStore, FilterRowPlugin, DefaultStringFilterOptions, AbsoluteHeightContainer } from "mosaic-data-table";
import { useContentManagerContext } from "./use-content-manager-context"
import { ActionList } from "@/shared/components/action-list";
import { ActionListItem } from "@/shared/components/action-list-item";
import { ConfirmationDialog } from "@/shared/components/confirmation-dialog";
import { CenterPagination, CoolPagination } from "@/shared/components/cool-pagination";
import { EmptyMessage } from "@/shared/components/empty-message";
import { MuiCardFooter } from "@/shared/components/mui-card-footer";
import { useDyamicGridFilter } from "./use-dyamic-grid-filter";
import Defaults from "@/config/Defaults.json";
import { useDynamicGridColumns } from "@/hooks/use-dynamic-grid-columns";
import { RelationViewerGridRoot } from "@/core-features/view-item/relation-viewer-grid";
import { useViewRelationStore } from "@/core-features/view-item/use-view-relation-store";
import { Edge } from "@/lib/apollo/graphql.entities";

export const ContentManagerScene = () => {

    const contentManager = useContentManagerContext();
    const contentManagerIds = useContentManagerIds(contentManager.contentManagerStore.state)
    const contentManagerSelection = useSelection<string>(contentManagerIds.storeIds);
    const [selectionEnabled, setSelectionEnabled] = useState(false);
    const [filterEnabled, setFilterEnabled] = useState<boolean>(!!Object.keys(contentManager.contentManagerSearch.state.filter).length);
    const [showId, setShowId] = useState<boolean>(false);
    const viewRelationStore = useViewRelationStore();
    const dynamicDialog = useDynamicDialog();
    const deleteConfirmationDialog = useDialog();
    const bulkDeleteConfirmationDialog = useDialog();

    useEffect(() => {
        viewRelationStore.closeAll();
    }, [contentManager.contentManagerStore.state.dataLoading]);

    const headCells = useDynamicGridColumns({
        entityName: contentManager.entityName,
        fields: contentManager.fullEntity?.fields,
        edges: contentManager.fullEntity?.edges,
        displayFieldName: contentManager.displayFieldName,
        //expansionStore: viewRelationStore,
		openRelation: useCallback((entityName: string, edge: Edge, entryId: string) => {
			viewRelationStore.addEdge(entryId, entityName, edge, entryId, true);
		}, [viewRelationStore]),
        showId: showId
    });

    const gridFilter = useDyamicGridFilter({
        fields: contentManager.fullEntity?.fields
    })

    const editEntry = useCallback((field: any) => {
        dynamicDialog.addPopup(ContentManagerEntryDialogContent, { entityName: contentManager.entityName, defaultValue: field }, FormOpenMode.Edit, field.id);
    }, [dynamicDialog, contentManager.entityName]);

    const deleteEntryConfirmed = useCallback(() => {
        if (!deleteConfirmationDialog.data) {
            return;
        }
        deleteConfirmationDialog.handleClose();
        contentManager.contentManagerStore.deleteItem(deleteConfirmationDialog.data as string);
    }, [deleteConfirmationDialog, contentManager.contentManagerStore.deleteItem]);

    const bulkDeleteEntryConfirmed = useCallback(() => {
        bulkDeleteConfirmationDialog.handleClose();
        contentManager.contentManagerStore.deleteBulk(contentManagerSelection.selected);
    }, [deleteConfirmationDialog, contentManager.contentManagerStore.deleteItem]);

    // Row Actions
    const actions: Action<unknown>[] = [
        {
            id: 'edit',
            render: (field: unknown) => (<MenuItem id='edit-menu-item' key={`edit-${field}`} onClick={() => editEntry(field)}>
                <ListItemIcon>
                    <EditIcon />
                </ListItemIcon>
                Edit
            </MenuItem>)
        },
        {
            id: 'remove',
            render: (field: any) => (<MenuItem id='remove-menu-item' key={`remove-${field}`} onClick={() => deleteConfirmationDialog.handleOpen(field.id)}> <ListItemIcon><DeleteIcon /></ListItemIcon> Remove </MenuItem>)
        },
    ];

    const finishDialog = useCallback((result: FinishResult) => {

        if (result.openMode == FormOpenMode.Edit) {
            if (!result.editId) {
                return;
            }
            contentManager.contentManagerStore.updateItem(result.editId, result.data);
            return;
        }

        if (result.openMode == FormOpenMode.New) {
            contentManager.contentManagerStore.addItem(result.data);
            return;
        }

    }, [contentManager.contentManagerStore]);

    const gridPlugins = useGridPlugins(
        CustomBodyCellContentRenderPlugin,

        usePluginWithParams(FilterRowPlugin, {
            visible: filterEnabled,
            filter: contentManager.contentManagerSearch.state.filter,
            filterChanged: contentManager.contentManagerSearch.handleFiltersApply,
            key: 'filter_row',
            filterColumns: gridFilter
        }),

        usePluginWithParams(PaddingPluggin, {}),
        usePluginWithParams(ColumnSortPlugin, {
            order: contentManager.contentManagerSearch.state.order,
            orderBy: contentManager.contentManagerSearch.state.sortBy,
            onSort: contentManager.contentManagerSearch.handleSortChange
        }),
        usePluginWithParams(RowSelectionPlugin, {
            visible: selectionEnabled,
            onGetRowId: contentManagerIds.getId,
            onSelectOne: contentManagerSelection.handleSelectOne,
            onDeselectOne: contentManagerSelection.handleDeselectOne,
            selectedIds: contentManagerSelection.selected
        }),
        usePluginWithParams(RowExpansionPlugin, {
            showExpanderButton: false,
            onGetRowId: contentManagerIds.getId,
            expanstionStore: viewRelationStore.expansionStore,
            getExpansionNode: useCallback((row: any, params: any) => (
                <AbsoluteHeightContainer>
                    <RelationViewerGridRoot rootEntryId={row.id} viewRelationStore={viewRelationStore} showId={showId}/>
                </AbsoluteHeightContainer>), [showId, viewRelationStore])
        }),
        ColumnsFillRowSpacePlugin,
        usePluginWithParams(RowActionsPlugin, {
            actions: actions
        }),
        usePluginWithParams(HighlightColumnPlugin, {}),
        PinnedColumnsPlugin,
        usePluginWithParams(SkeletonLoadingPlugin, {
            isLoading: contentManager.contentManagerStore.state.dataLoading,
            rowsWhenEmpty: Defaults.dataTable.skeletonRowsCount,
            maxRowsWhenNotEmpty: 15
        }),
        usePluginWithParams(EmptyDataPlugin, {
            content: <EmptyMessage />
        }),
    );

    return (
        <>
            <Card>
                <CardHeader title={<ContentManagerSearch
                    headCells={headCells}
                    disabled={false}
                    onRefresh={contentManager.contentManagerStore.refresh}
                    filters={contentManager.contentManagerSearch.state.advancedFilters}
                    onFiltersApply={contentManager.contentManagerSearch.handleAdvancedFiltersApply}
                    onFiltersClear={contentManager.contentManagerSearch.handleAdvancedFiltersClear}
                    onBulkDelete={() => bulkDeleteConfirmationDialog.handleOpen()}
                    onQueryChange={contentManager.contentManagerSearch.handleQueryChange}
                    query={contentManager.contentManagerSearch.state.query}
                    selected={contentManagerSelection.selected}

                    selectionEnabled={selectionEnabled}
                    onSelectionEnabledChange={(value) => {
                        setSelectionEnabled(value)
                        if (!value) {
                            contentManagerSelection.handleDeselectAll();
                        }
                    }}

                    filterEnabled={filterEnabled}
                    onFilterEnabledChange={setFilterEnabled}

                    contentManagerSearch={contentManager.contentManagerSearch}
                    showId={showId}
                    onShowIdChange={setShowId}
                />
                }>

                </CardHeader>

                <Divider />

                <CardContent>


                    <MosaicDataTable
                        plugins={gridPlugins}
                        caption={`${contentManager.entityName} content`}
                        items={contentManager.contentManagerStore.state.data}
                        headCells={headCells}
                    />

                    <CenterPagination>
                        <CoolPagination
                            page={contentManager.contentManagerSearch.state.page}
                            pagesCount={contentManager.contentManagerStore.state.pagesCount}
                            totalItems={contentManager.contentManagerStore.state.totalItems}
                            onChange={contentManager.contentManagerSearch.handlePageChange}
                        />
                    </CenterPagination>
                </CardContent>

                <Divider />

                <MuiCardFooter>
                    <ActionList>
                    <ActionListItem
                            onClick={() => {
                                dynamicDialog.addPopup(ContentManagerEntryDialogContent, { entityName: contentManager.entityName }, FormOpenMode.New);
                            }}
                            icon={(
                                <SvgIcon fontSize="small">
                                    <PlusCircleIcon />
                                </SvgIcon>
                            )}
                            label="Add New"
                            aria-label="Add new entry"
                            aria-haspopup="dialog" />
                            
                    </ActionList>
                </MuiCardFooter>

            </Card>

            <DynamicDialog
                ref={dynamicDialog.ref}
                finish={finishDialog} />

            <ConfirmationDialog
                message="Are you sure you want to delete this item? This can't be undone."
                onCancel={deleteConfirmationDialog.handleClose}
                onConfirm={deleteEntryConfirmed}
                open={deleteConfirmationDialog.open}
                title="Delete Confirmation"
                variant="error"
            />

            <ConfirmationDialog
                message="Are you sure you want to delete the selected items? This can't be undone."
                onCancel={bulkDeleteConfirmationDialog.handleClose}
                onConfirm={bulkDeleteEntryConfirmed}
                open={bulkDeleteConfirmationDialog.open}
                title="Delete Confirmation"
                variant="error"
            />
        </>
    )
}
ContentManagerScene.displayName = 'ContentManagerScene';