"use client";

import { useEntities } from "@/hooks/use-entities";
import { CategoryCollapse } from "@/shared/components/category-collapse";
import { TreeList, TreeListItemNavigation, TreeListItemSkeleton, TreeListItemText } from "@/shared/components/tree-list";
import { useSkipLink } from "@/shared/layouts/header/skip-links/src/skip-links";
import { NameCaptionEntity } from "@/types/entity";
import { Box, Chip, CircularProgress, Stack, Tooltip } from "@mui/material";
import { useMemo } from "react";
import Defaults from "@/config/Defaults.json";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { MenuItem } from "@/shared/layouts/header/Topnav/top-nav";
import ArrowRight from "@/assets/icons/iconly/bulk/arrow-right";
import { GRAPHQL_ENTITY_PLAYGROUND_URL, GRAPHQL_QUERY_PLAYGROUND_URL } from "@/config/CONST";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import { HiExternalLink } from "react-icons/hi";

export default function DeveloperEntities() {

    const skipId = useSkipLink({ id: 'developer-entities', title: 'Go To Entities' });

    const { entities, loadingNameCaptionEntities: loading } = useEntities();
    const entityId = useCurrentEntityNameContext();

    return (
        <Box id={skipId}>

            <CategoryCollapse
                label={<Stack direction="row" alignItems="center" gap={1}>
                    USEFULL
                </Stack>}
                openImmediately={true}>
                <TreeList>

                    <TreeListItemNavigation
                        active={false}
                        path={`/developer/graphql/export`}
                        label="Export Collection" />

                    <TreeListItemNavigation
                        active={false}
                        path={GRAPHQL_QUERY_PLAYGROUND_URL!}
                        label="Playground"
                        externalLink={true}
                        icon={<HiExternalLink size="14px"/>} />

                    <TreeListItemNavigation
                        active={false}
                        path={GRAPHQL_ENTITY_PLAYGROUND_URL!}
                        label="Entity Playground"
                        externalLink={true}
                        icon={<HiExternalLink size="14px"/>} />

                </TreeList>
            </CategoryCollapse>

            <CategoryCollapse
                label={<Stack direction="row" alignItems="center" gap={1}>
                    GRAPH-QL API
                </Stack>}
                openImmediately={true}>


                <TreeList>
                    <TreeListItemNavigation
                        active={false}
                        path={`/developer/graphql`}
                        label="Entities" />

                    {entities.map((entity, i) => (
                        <EntityButtonNavigation
                            key={entity.name}
                            active={entityId == entity.name}
                            entity={entity} />))
                    }

                </TreeList>
            </CategoryCollapse>
        </Box>
    )
}



interface EntityButtonNavigationProps {
    entity: NameCaptionEntity,
    active: boolean
}
const EntityButtonNavigation = (props: EntityButtonNavigationProps) => {
    const { entity: { name: entityName, caption: entityCaption }, active } = props;

    const appStatus = useAppStatus();

    const busyIcon = useMemo(() => {
        return (<Tooltip title="Code generation" placement='bottom' arrow>
            <CircularProgress size={15} />
        </Tooltip>);

    }, []);

    if (appStatus.isEntityInCodeGeneration(entityName)) {
        return (<TreeListItemText
            slotProps={{ primary: { color: 'text.secondary' } }}
            label={entityCaption ?? entityName}
            icon={busyIcon} />)
    }

    return (<TreeListItemNavigation
        active={active}
        path={`/developer/graphql/${entityName}`}
        label={entityCaption ?? entityName} />)
}
