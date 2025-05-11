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
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";

export default function ContentManagernEntities() {

    const skipId = useSkipLink({ id: 'entity-manager-entities', title: 'Go To Entities' });

    const entityId = useCurrentEntityNameContext();
    const { entities, loadingNameCaptionEntities: loading } = useEntities();

    return (
        <Box id={skipId}>
            <CategoryCollapse
                label={<Stack direction="row" alignItems="center" gap={1}>
                    ENTITIES
                    {!loading && <Chip component="span" label={entities.length} color="secondary" variant="outlined" size="small" />}
                </Stack>}
                openImmediately={true}>
                <TreeList>

                    {!loading && entities.map((entity, i) => (
                        <EntityButtonNavigation
                            key={entity.name}
                            active={entityId == entity.name}
                            entity={entity} />))
                    }

                    {loading && Array(entities.length || Defaults.entityList.skeletonRowsCount).fill(0).map((_, i) => (
                        <TreeListItemSkeleton key={i} />
                    ))}

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
        path={`/content-manager/${entityName}`}
        label={entityCaption ?? entityName} />)
}
