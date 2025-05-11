import { Edge } from "@/lib/apollo/graphql.entities";
import { useMemo } from "react";
import { ManyRelationViewerGrid } from "./many-relation-viewer-grid";
import { OneRelationViewerGrid } from "./one-relation-viewer-grid";
import { Box } from "@mui/material";

interface RelationViewerGridRootProps {
    entityName: string;
    entryId: string;
    showId: boolean;
    edge: Edge;
}
export const RelationViewerGridRoot = (props: RelationViewerGridRootProps) => {
    return (
        <Box sx={{
            paddingBottom: '10px',
            paddingX: '10px',
            backgroundImage: 'url(/rough-diagonal.png)',
        }}>
            <RelationViewerGrid {...props} />
        </Box>
    )
}



export const RelationViewerGrid = (props: RelationViewerGridRootProps) => {

    const isOne = useMemo(() => {
        return props.edge.relationType === 'One' || props.edge.relationType === 'OneToMany' || props.edge.relationType === 'OneToOne';
    }, [props.edge.relationType]);

    if (isOne) {
        return <OneRelationViewerGrid {...props} />;
    }

    return <ManyRelationViewerGrid {...props} />;
}
