import { useCallback, useMemo } from "react";
import { useOneViewRelationStore, useViewRelationStore } from "./use-view-relation-store";
import { useContentManagerSearch } from "@/hooks/use-content-manager-search";
import { useContentManagerStore } from "@/hooks/use-content-manager-store";
import { useDynamicGridColumns } from "@/hooks/use-dynamic-grid-columns";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { Action, ColumnsFillRowSpacePlugin, CustomBodyCellContentRenderPlugin, EmptyDataPlugin, Filter, HighlightColumnPlugin, MosaicDataTable, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, useGridPlugins, usePluginWithParams } from "mosaic-data-table";
import { eqStringFoldOperator } from "../dynamic-filter/filter-operators";
import { Box, Button, IconButton, Link, ListItemIcon, MenuItem, Stack, styled, Typography } from "@mui/material";
import DirectionsRunIcon from '@mui/icons-material/DirectionsRun';
import { useRouter } from "next/navigation";
import { Edge } from "@/lib/apollo/graphql.entities";
import { getAdvancedFiltersFromGridFilter } from "@/lib/utils/get-advanced-filters-from-grid-filters";
import { ArrowRoot, RoundPanelPlaceholder } from "./style";
import CloseIcon from '@mui/icons-material/Close';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import NextLink from 'next/link';
import { EmptyMessage } from "@/shared/components/empty-message";
import { useGetEdgeValue } from "@/hooks/use-get-edge-value";

interface RelationViewerGridRootProps {
	showId: boolean
	rootEntryId: string;
	viewRelationStore: ReturnType<typeof useViewRelationStore>
}
export const RelationViewerGridRoot = (props: RelationViewerGridRootProps) => {

	const { push } = useRouter();
	const oneViewRelationStore = useOneViewRelationStore(props.rootEntryId, props.viewRelationStore);
	const { visibleEdge } = oneViewRelationStore;

	const hasHistory = useMemo(() => oneViewRelationStore.edges.length > 1, [oneViewRelationStore.edges]);

	const { entities: nameCaptionEntities } = useEntities();
	const nameCaptionEntity = useMemo(() => nameCaptionEntities.find(i => i.name == visibleEdge?.entityName), [visibleEdge?.entityName, nameCaptionEntities]);

	if (!visibleEdge) {
		return;
	}

	return (
		<Box sx={{
			padding: '10px',
			backgroundImage: 'url(/rough-diagonal.png)',
		}}>

			<Box sx={{
				position: 'relative',
			}}>

				<Stack gap={1}>

					<RoundPanelPlaceholder sx={{
						backdropFilter: 'blur(100px)'
					}}>
						<Stack direction="row" gap={1} alignItems="center" justifyContent="space-between">
							<Stack direction="row" gap={1} alignItems="center">

								<IconButton aria-label="back" size="medium" onClick={oneViewRelationStore.goBack}>
									<ArrowBackIcon fontSize="inherit" />
								</IconButton>

								{hasHistory && (<Typography variant="h6">
									{'.'.repeat(oneViewRelationStore.edges.length)} /
								</Typography>)}

								<Stack direction="row" gap={1} alignItems="center">

									<Typography variant="h6">
										{nameCaptionEntity?.caption}
									</Typography>
									➝
									<Typography variant="h6" sx={{ fontStyle: 'italic' }}>
										{visibleEdge.edge.caption}
									</Typography>
									➝
									<Link component={NextLink} href={`/content-manager/${visibleEdge?.edge.relatedEntity.name}`} color="inherit" >
										<Typography variant="h6">
											{visibleEdge.edge.relatedEntity.caption}
										</Typography>
									</Link>

								</Stack>

							</Stack>

							<IconButton aria-label="close" size="medium" onClick={oneViewRelationStore.close}>
								<CloseIcon fontSize="inherit" />
							</IconButton>

						</Stack>
					</RoundPanelPlaceholder>

					<RoundPanelPlaceholder>
						<RelationViewerGrid key={`${visibleEdge.entityName}-${visibleEdge.entryId}`} entityName={visibleEdge.entityName} entryId={visibleEdge.entryId} edge={visibleEdge.edge} showId={props.showId} openRelation={oneViewRelationStore.addEdge} />
					</RoundPanelPlaceholder>
				</Stack>
			</Box>
		</Box>
	)
}




interface RelationViewerGridProps {
	entityName: string;
	entryId: string;
	edge: Edge;
	showId: boolean;

	openRelation: (entityName: string, edge: Edge, entryId: string) => void;
}
const RelationViewerGrid = (props: RelationViewerGridProps) => {

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

	const fullEntity = useFullEntity({ entityName: props.edge.relatedEntity.name });
	
	const relationData = useGetEdgeValue({
		entityName: props.entityName,
		entryId: props.entryId,
		edge: props.edge,
		fields: 'allfields',
		edges: useMemo(() => ([
			...fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
				.filter(i => i.relationType !== 'ManyToOne')
				.filter(i => i.relationType !== 'Many')
				.map(i => ({
					name: i.name,
					fields: ['id', i.relatedEntity.displayField.name],
				})) ?? []
		]), [fullEntity?.edges])
	})

	const gridData = useMemo(() => {
		if(!relationData.data){
			return [];
		}
		if(Array.isArray(relationData.data)){
			return relationData.data;
		}
		return [relationData.data];
	}, [relationData.data]);


	const headCells = useDynamicGridColumns({
		entityName: fullEntity?.name ?? '',
		fields: fullEntity?.fields,
		edges: fullEntity?.edges,
		displayFieldName: fullEntity?.displayField?.name ?? 'id',
		openRelation: props.openRelation,
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
		ColumnsFillRowSpacePlugin,
		usePluginWithParams(HighlightColumnPlugin, {}),

		usePluginWithParams(RowActionsPlugin, {
			actions: actions
		}),
		usePluginWithParams(EmptyDataPlugin, {
			content: <EmptyMessage />
		}),
		PinnedColumnsPlugin
	)

	return (<MosaicDataTable
		plugins={gridPlugins}
		caption={`${props.entityName} Entry Viewer`}
		items={gridData}
		headCells={headCells}
	/>)
}