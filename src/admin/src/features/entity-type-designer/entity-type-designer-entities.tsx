'use client'

import { Box, SvgIcon, CircularProgress, Tooltip, Stack, Chip } from "@mui/material";
import { useParams, useRouter } from "next/navigation"
import { useCallback, useMemo } from "react";
import PlusIcon from "@heroicons/react/24/outline/PlusIcon";
import PencilIcon from "@heroicons/react/24/outline/PencilIcon";
import TrashIcon from "@heroicons/react/24/outline/TrashIcon";
import { EntityTypeNewEntityDialog } from "./entity-type-new-entity-dialog";
import { DynamicDialog } from "../../core-features/dynamic-dialog/src/dynamic-dialog";
import { FormOpenMode } from "@/core-features/dynamic-form/form-field";
import { useDynamicDialog } from "@/core-features/dynamic-dialog/src/use-dynamic-dialog";
import { Theme, styled } from '@mui/material/styles';
import { useEntitiesDesigner } from "./use-designer-entities";
import { useSkipLink } from "@/shared/layouts/header/skip-links/src/skip-links";
import { CategoryCollapse } from "@/shared/components/category-collapse";
import { TreeList, TreeListItemButton, TreeListItemNavigation, TreeListItemSkeleton, TreeListItemText } from "@/shared/components/tree-list";
import { useEntities } from "@/hooks/use-entities";
import { ChangedNameCaptionEntity } from "@/types/entity";
import Defaults from "@/config/Defaults.json";
import { useAppStatus } from "@/store/app-state/use-app-state";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";

const iconStatus = {
    'new': <SvgIcon fontSize="small"><PlusIcon /></SvgIcon>,
    'edited': <SvgIcon fontSize="small"><PencilIcon /></SvgIcon>,
    'deleted': <SvgIcon fontSize="small"><TrashIcon /></SvgIcon>,
    'unchanged': null
}

const newEntityBorderStyle = (theme: Theme) => {
    return `2px dashed ${theme.palette.mode == 'dark' ? theme.palette!.background!.paper : theme.palette!.neutral[100]}`;
}

const TreeListItemButtonStyled = styled(TreeListItemButton)(({ theme }) => [
    {
        color: 'text.secondary',
        borderWidth: '2px',
        borderStyle: 'dashed'
    },
    theme.applyStyles('light', {
        borderColor: 'var(--mui-palette-neutral-100)',

    }),
    theme.applyStyles('dark', {
        borderColor: 'var(--mui-palette-background-paper)',
    }),
]);


export default function EntityTypeBuilderEntities() {

    const entityId = useCurrentEntityNameContext();
    const router = useRouter();
    const { allNameCaptionEntities, addEntity } = useEntitiesDesigner();
    const graphEntities = useEntities();
    const dynamicDialog = useDynamicDialog();

    const skipId = useSkipLink({ id: 'entity-manager-entities', title: 'Go To Entities' });

    const finish = useCallback(({ data, openMode }: {
        data: ChangedNameCaptionEntity,
        openMode?: FormOpenMode
    }) => {
        addEntity(data);
        router.push(`/entity-type-designer/${data.name}`);

    }, [addEntity]);

    return (
        <Box id={skipId}>

            <CategoryCollapse label={<Stack direction="row" alignItems="center" gap={1}>
                ENTITIES
                {!graphEntities.loadingNameCaptionEntities && <Chip component="span" label={allNameCaptionEntities.length} color="secondary" variant="outlined" size="small" />}
            </Stack>}
                openImmediately={true}>
                <TreeList>

                    {!graphEntities.loadingNameCaptionEntities && allNameCaptionEntities.map((entity, i) => (
                        <EntityButtonNavigation
                            key={entity.name}
                            active={entityId == entity.name}
                            entity={entity} />))
                    }

                    {graphEntities.loadingNameCaptionEntities && Array(allNameCaptionEntities.length || Defaults.entityList.skeletonRowsCount).fill(0).map((_, i) => (
                        <TreeListItemSkeleton key={i} />
                    ))}

                    <TreeListItemButtonStyled
                        label="Add Entity"
                        icon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                        onClick={() => dynamicDialog.addPopup(EntityTypeNewEntityDialog, {}, FormOpenMode.New)}
                        aria-label="Add new entity"
                        aria-haspopup="dialog"
                    />
                </TreeList>
            </CategoryCollapse>

            <DynamicDialog
                ref={dynamicDialog.ref}
                finish={finish} />
        </Box>
    )
}

interface EntityButtonNavigationProps {
    entity: ChangedNameCaptionEntity,
    active: boolean
}
const EntityButtonNavigation = (props: EntityButtonNavigationProps) => {
    const { entity: { name: entityName, caption: entityCaption, designerStatus }, active } = props;
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
        path={`/entity-type-designer/${entityName}`}
        label={entityCaption ?? entityName}
        icon={iconStatus[designerStatus]} />)
}
