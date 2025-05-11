import { useContentManagerSearch } from "@/hooks/use-content-manager-search";
import { useContentManagerStore } from "@/hooks/use-content-manager-store";
import { useDynamicGridColumns } from "@/hooks/use-dynamic-grid-columns";
import { useFullEntity } from "@/hooks/use-entities";
import { AbsoluteHeightContainer, Action, ColumnsFillRowSpacePlugin, CustomBodyCellContentRenderPlugin, Filter, FilterRowPlugin, HighlightColumnPlugin, MosaicDataTable, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, RowExpansionPlugin, useGridPlugins, usePluginWithParams, useRowExpansionStore } from "mosaic-data-table";
import { eqStringFoldOperator } from "../dynamic-filter/filter-operators";
import { useContentManagerIds } from "@/features/content-manager/use-content-manger-ids";
import { useCallback, useMemo } from "react";
import { Box, ListItemIcon, MenuItem, Stack, styled } from "@mui/material";
import DirectionsRunIcon from '@mui/icons-material/DirectionsRun';
import { useRouter } from "next/navigation";
import { Edge } from "@/lib/apollo/graphql.entities";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { ArrowRoot } from "./style";
import { RelationViewerGridRoot } from "./relation-viewer-grid";

interface OneItemViewerGridProps {
    entityName: string;
    entryId: string;
    edge: Edge;
    showId: boolean;
}
export const ManyRelationViewerGrid = (props: OneItemViewerGridProps) => {

    const { push } = useRouter();
    
    /*
    entityName
    ├─ ☐ field_01
    ├─ ☐ field_02
    ├─ ☐ field_03
    └─ ☐ EDGE_01
        ├─ ☑ field_01_01
        ├─ ☑ field_01_02
        ├─ ☑ field_01_03
        └─ ☑ EDGE_01_01
            ├─ ☑ id
            └─ ☑ [EDGE_01_01.relatedEntity.displayField.name]
    */

    const contentManagerSearch = useContentManagerSearch({
        initialFilter: {
            id: {
                operator: eqStringFoldOperator.name,
                value: props.entryId
            }
        }
    });

    const fullEntity = useFullEntity({ entityName: props.edge.relatedEntity.name });
    const selectFields = useMemo(() => fullEntity?.fields.map(i => i.name) ?? [], [fullEntity?.fields]);
    const selectEdges = useMemo(() => ([
        ...fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
            .filter(i => i.relationType !== 'ManyToOne')
            .filter(i => i.relationType !== 'Many')
            .map(i => ({
                name: i.name,
                fields: ['id', i.relatedEntity.displayField.name],
            })) ?? []
    ]), [fullEntity?.edges]);

    const contentManagerStore = useContentManagerStore({
        entityName: props.entityName,

        page: contentManagerSearch.state.page,
        rowsPerPage: contentManagerSearch.state.rowsPerPage,
        sortBy: contentManagerSearch.state.sortBy,
        order: contentManagerSearch.state.order,
        edges: useMemo(() => {
            if (!selectFields.length && !selectEdges.length) {
                return undefined;
            }

            return [{
                name: props.edge.name,
                fields: selectFields,
                edges: selectEdges
            }]
        }, [selectFields, selectEdges]),

        filters: useMemo(() => getAdvancedFiltersFromGridFilter(contentManagerSearch.state.filter), [contentManagerSearch.state.filter]),
    });


    const gridData = useMemo(() => {
        if(!contentManagerStore.state.data?.length){
            return []
        }

        return contentManagerStore.state.data[0][props.edge.name];
    }, [contentManagerStore.state.data]);

    const expansionStore = useRowExpansionStore();
    
    const headCells = useDynamicGridColumns({
        entityName: fullEntity?.name ?? '',
        fields: fullEntity?.fields,
        edges: fullEntity?.edges,
        displayFieldName: fullEntity?.displayField?.name ?? 'id',
        expansionStore: expansionStore,
        showId: props.showId
    });

    // Row Actions
    const actions: Action<unknown>[] = [
        {
            id: 'goto',
            render: (field: unknown) => (<MenuItem id='edit-menu-item' key={`edit-${field}`} onClick={() => gotoEntry(field)}>
                <ListItemIcon>
                    <DirectionsRunIcon />
                </ListItemIcon>
                Go To Entry
            </MenuItem>)
        },
    ];

    const gotoEntry = useCallback((entry: any) => {
        const path = `/content-manager/${props.edge.relatedEntity.name}`;

        const filter: Filter = {
            "id": {
                operator: eqStringFoldOperator.name,
                value: entry.id
            }
        }

        const withFilters = `${path}?f=${JSON.stringify(filter)}`;
        push(withFilters);
    }, []);

    const gridPlugins = useGridPlugins(
        CustomBodyCellContentRenderPlugin,
        usePluginWithParams(PaddingPluggin, {}),
        usePluginWithParams(RowExpansionPlugin, {
            showExpanderButton: false,
            onGetRowId: useCallback((i: any) => i.id, []),
            expanstionStore: expansionStore,
            getExpansionNode: useCallback((row: any, params: any) => (
                <AbsoluteHeightContainer sx={{
                    border: '1px solid var(--mui-palette-background-paper)',
                }}>
                    <RelationViewerGridRoot entityName={params.entityName} edge={params.edge} entryId={params.entryId} showId={props.showId}/>
                </AbsoluteHeightContainer>), [props.showId])
        }),
        ColumnsFillRowSpacePlugin,
        usePluginWithParams(HighlightColumnPlugin, {}),

        usePluginWithParams(RowActionsPlugin, {
            actions: actions
        }),
        PinnedColumnsPlugin
    )

    return (
        <Stack gap={0}>

            <Box sx={{
                backgroundColor: 'var(--mui-palette-background-paper)',
                borderRadius: 2,
                marginTop: '10px',
                position: 'relative',
            }}>
                <ArrowRoot />
                <Box sx={{
                    borderRadius: 2,
                    overflow: 'hidden',
                }}>
                    <MosaicDataTable
                        plugins={gridPlugins}
                        caption={`${props.entityName} Entry Viewer`}
                        items={gridData}
                        headCells={headCells}
                    />
                </Box>
            </Box>
        </Stack>)
}


interface OneItemViewerGridRootProps {
    entityName: string;
    entryId: string;
    showId: boolean;
    edge: Edge;
}
export const ManyRelationViewerGridRoot = (props: OneItemViewerGridRootProps) => {

    return (

        <Box sx={{
            paddingBottom: '10px',
            paddingX: '10px',
            backgroundImage: 'url(/rough-diagonal.png)',
        }}>
            <ManyRelationViewerGrid
                key={`${props.entityName}${props.entryId}`}
                entityName={props.entityName}
                entryId={props.entryId}
                edge={props.edge}
                showId={props.showId} />
        </Box>
    )
}
