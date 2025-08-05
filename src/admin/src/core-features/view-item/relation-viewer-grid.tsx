import { useCallback, useEffect, useMemo, useRef } from "react";
import { useOneViewRelationStore, useViewRelationStore } from "./use-view-relation-store";
import { useDynamicGridColumns } from "@/hooks/use-dynamic-grid-columns";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import { Action, ColumnsFillRowSpacePlugin, CustomBodyCellContentRenderPlugin, EmptyDataPlugin, Filter, HighlightColumnPlugin, MosaicDataTable, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, useGridPlugins, usePluginWithParams } from "mosaic-data-table";
import { eqStringFoldOperator } from "../dynamic-filter/filter-operators";
import { Box, Button, IconButton, Link, ListItemIcon, MenuItem, Stack, styled, Tooltip, Typography } from "@mui/material";
import DirectionsRunIcon from '@mui/icons-material/DirectionsRun';
import { useRouter } from "next/navigation";
import { Edge } from "@/lib/apollo/graphql.entities";
import { ArrowRoot, RoundPanelPlaceholder } from "./style";
import CloseIcon from '@mui/icons-material/Close';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import NextLink from 'next/link';
import { EmptyMessage } from "@/shared/components/empty-message";
import { useRelationContentManagerStore } from "@/hooks/use-relation-content-manager-store";
import { Key } from "@/shared/components/key-handler/types";
import { ShortcutIconButton } from "@/shared/components/key-handler/with-click-shortcut";
import { useAutoFocusFirstElementOnce } from "@/hooks/use-auto-focus-first-element-once";

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

	const autoFocusHandler = useAutoFocusFirstElementOnce()
	const hostRef = useRef<HTMLDivElement | null>(null);

	const setHostRef = useCallback((instance: HTMLDivElement | null) => {
		autoFocusHandler.setRef(instance);
		hostRef.current = instance;
	}, [autoFocusHandler.setRef]);

	const goBack = useCallback(() => {
		oneViewRelationStore.goBack();
		autoFocusHandler.focus();
	}, [oneViewRelationStore.goBack]);

	const close = useCallback(() => {
		oneViewRelationStore.close();
	}, [oneViewRelationStore.close]);


	useEffect(() => {
		if(hostRef.current){
			hostRef.current.focus();
		}
	}, [`${visibleEdge?.entityName}-${visibleEdge?.entryId}`]);
	
	if (!visibleEdge) {
		return;
	}

	return (
		<Box
			ref={setHostRef} 
			tabIndex={0}
			sx={{
				padding: '10px',
				// backgroundImage: 'url(/rough-diagonal.png)',
				//background: 'url(/assets/img/relation-background.svg)',
				backgroundColor: 'var(--mui-palette-background-default)',
				'&:focus': {
					outline: 'none'
				}
			}}>

			<Box sx={{
				position: 'relative',
			}}>

				<Stack gap={1}>

					<RoundPanelPlaceholder sx={{
						// backdropFilter: 'blur(100px)'
					}}>
						<Stack direction="row" gap={1} alignItems="center" justifyContent="space-between">
							<Stack direction="row" gap={1} alignItems="center">

								<ShortcutIconButton 
									shortcutKey={Key.b}
									shortcutTarget={hostRef}
									tooltip="Back"
									aria-label="back" size="medium" onClick={goBack}>
									<ArrowBackIcon fontSize="inherit" />
								</ShortcutIconButton>
								

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

							<ShortcutIconButton 
								shortcutKey={Key.Escape}
								shortcutTarget={hostRef}
								tooltip="Back"
								aria-label="close" size="medium" onClick={close}>
								<CloseIcon fontSize="inherit" />
							</ShortcutIconButton>

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

	const relationData = useRelationContentManagerStore({
		entityName: props.entityName,
		entryId: props.entryId,
		edge: props.edge,
		fields: 'grid'
	});
	// const relationData = useGetEdgeValue({
	// 	entityName: props.entityName,
	// 	entryId: props.entryId,
	// 	edge: props.edge,
	// 	fields: 'allfields',
	// 	edges: useMemo(() => ([
	// 		...fullEntity?.edges.filter(i => i.relationType !== 'ManyToMany')
	// 			.filter(i => i.relationType !== 'ManyToOne')
	// 			.filter(i => i.relationType !== 'Many')
	// 			.map(i => ({
	// 				name: i.name,
	// 				fields: ['id', i.relatedEntity.displayField.name],
	// 			})) ?? []
	// 	]), [fullEntity?.edges])
	//})

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