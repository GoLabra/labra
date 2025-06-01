import { Edge } from "@/lib/apollo/graphql.entities";
import { useRowExpansionStore } from "mosaic-data-table";
import { useCallback, useMemo } from "react";

export type EdgeNode = {
	entityName: string;
	edge: Edge;
	entryId: string;
}
export type RelationViewState = {
	edges: EdgeNode[];
}

export const useViewRelationStore = () => {
	const expansionStore = useRowExpansionStore<RelationViewState>();

	const close = useCallback((rootEntryId: string) => {
		expansionStore.setParams({
			rowId: rootEntryId,
			params: {
				edges: []
			},
			openImmediately: false
		});
	}, [expansionStore]);

	const closeAll = useCallback(() => {
		expansionStore.clear();
	}, [expansionStore]);

	const addEdge = useCallback((rootEntryId: string, entityName: string, edge: Edge, entryId: string, addAsFirst: boolean = false) => {

		if (addAsFirst) {
			expansionStore.setParams({
				rowId: rootEntryId,
				params: {
					edges: [
						{
							entityName,
							edge,
							entryId
						}]
				},
				openImmediately: true
			});

			return;
		}

		const edgeNodesInfo = expansionStore.getExpansionInfo(rootEntryId);
		const params = edgeNodesInfo.params ?? {};
		const edges = params?.edges ?? [];

		expansionStore.setParams({
			rowId: rootEntryId,
			params: {
				edges: [
					...edges,
					{
						entityName,
						edge,
						entryId
					}]
			},
			openImmediately: true
		});

	}, [expansionStore]);

	const goBack = useCallback((rootEntryId: string) => {
		const edgeNodesInfo = expansionStore.getExpansionInfo(rootEntryId);
		const params = edgeNodesInfo.params ?? {};
		const edges = params?.edges ?? [];

		expansionStore.setParams({
			rowId: rootEntryId,
			params: {
				edges: [
					...edges.slice(0, -1)
				]
			},
			openImmediately: true
		});
	}, [expansionStore]);

	return useMemo(() => ({
		expansionStore,
		addEdge,
		goBack,
		close,
		closeAll
	}), [expansionStore, goBack, addEdge, close, closeAll]);
}


export const useOneViewRelationStore = (rootEntryId: string, viewRelationStore: ReturnType<typeof useViewRelationStore>) => {

	const edges = useMemo(() => {
		const edgeNodesInfo = viewRelationStore.expansionStore.getExpansionInfo(rootEntryId);
		const params = edgeNodesInfo.params ?? {};
		const edges = params?.edges ?? [];
		return edges;
	}, [rootEntryId, viewRelationStore.expansionStore.getExpansionInfo]);

	const visibleEdge = useMemo(() => {
		return edges.findLast(i => i);
	}, [edges]);

	const addEdge = useCallback((entityName: string, edge: Edge, entryId: string) => {
		viewRelationStore.addEdge(rootEntryId, entityName, edge, entryId); 
	}, [viewRelationStore, rootEntryId]);

	const goBack = useCallback(() => {
		viewRelationStore.goBack(rootEntryId);
	}, [viewRelationStore, rootEntryId]);

	const close = useCallback(() => {
		viewRelationStore.close(rootEntryId);
	}, [viewRelationStore, rootEntryId]);

	return useMemo(() => ({
		edges,
		visibleEdge,
		addEdge,
		goBack,
		close
	}), [edges, visibleEdge, addEdge, goBack, close]);
}