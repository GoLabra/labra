import { Box, Divider, Grid, IconButton, ListItemIcon, ListItemText, MenuItem, Paper, Stack, SvgIcon, Typography, } from "@mui/material";
import { useCallback, useEffect, useMemo } from "react";
import { getChildTypeForEntityChild, getIconForEntityChild } from "./designer-field-map";
import HistoryIcon from "@mui/icons-material/History";
import EditIcon from "@mui/icons-material/Edit";
import AddIcon from "@mui/icons-material/Add";
import DeleteIcon from "@mui/icons-material/Delete";
import EntityFieldType from "./entity-field-type-descriptor";
import { DesignerEdge, DesignerEntityStatus, DesignerField } from "@/types/entity";
import { alpha } from "@mui/material/styles";
import { EntityTypeDesignerFieldFlags } from "./entity-type-designer-field-flags";
import { useEntitiesDesigner, useEntityDesignerForEntity } from "./use-designer-entities";
import { ColumnsFillRowSpacePlugin, CustomBodyCellContentRenderPlugin, EmptyDataPlugin, ColumnDef, HideHeaderPlugin, MosaicDataTable, PaddingPluggin, PinnedColumnsPlugin, RowActionsPlugin, useGridPlugins, usePluginWithParams } from "mosaic-data-table";
import { Stats } from "@/shared/components/stats";
import { EmptyMessage } from "@/shared/components/empty-message";
import { ActionsButton } from "@/shared/components/actions-button";
import Avvvatars from "avvvatars-react";
import { MenuButton } from "@/shared/components/menu/menu-button";
import MoreVertIcon from '@mui/icons-material/MoreVert';

const EntityChangesIcon = ({ designerStatus }: { designerStatus: DesignerEntityStatus }) => {

    const icon = useMemo(() => {
        switch (designerStatus) {
            case 'new':
                return <AddIcon />;
            case 'edited':
                return <EditIcon />;
            case 'deleted':
                return <DeleteIcon />;
        }
    }, [designerStatus]);

    return (
        <Box
            sx={{
                borderWidth: "1px",
                borderStyle: "dashed",
                borderColor: (theme) => alpha(theme.palette.neutral[400], 0.7),
                color: (theme) => alpha(theme.palette.neutral[400], 0.7),
                padding: 0.4,
                display: "flex",
                borderRadius: 1,
                width: "32px",
                height: "32px",

            }}>
            {icon}
        </Box>
    );
};

interface EntityTypeChangesProps {
    entityName: string;
    onClose: () => void;
}
export const EntityTypeChanges = (props: EntityTypeChangesProps) => {

    const entityDesignerChanges = useEntityDesignerForEntity(props.entityName);

    const allFieldsCount = useMemo((): number => {        
        return entityDesignerChanges.justChanges?.fields.length ?? 0;
    }, [entityDesignerChanges.justChanges?.fields]);

    const allEdgesCount = useMemo((): number => {
        return entityDesignerChanges.justChanges?.edges.length ?? 0;
    }, [entityDesignerChanges.justChanges?.edges]);

    useEffect(() => {
        if(entityDesignerChanges.hasChanges == false){
            props.onClose();
        }
    }, [entityDesignerChanges.hasChanges]);

    if(entityDesignerChanges.hasChanges == false){
        return <Typography>No changes!</Typography>;
    }

    return (
        <>
            <Box>
                <Grid container spacing={1}>
                    <Stats
                        label="Entities"
                        value="1"
                    ></Stats>
                    <Stats
                        label="Fields"
                        value={allFieldsCount.toString()}
                    ></Stats>
                    <Stats
                        label="Edges"
                        value={allEdgesCount.toString()}
                    ></Stats>
                </Grid>
            </Box>

            <ChangesGrid entityName={props.entityName} onClose={props.onClose}></ChangesGrid>
        </>
    );
};

// Header
const headerCells: ColumnDef<DesignerField | DesignerEdge>[] = [
    {
        id: "designer-status",
        header: '',
        cell: (child: DesignerField | DesignerEdge) => <EntityChangesIcon key={child.name} designerStatus={child.designerStatus} />,
        width: 40,
    },
    {
        id: "name",
        header: "Type",
        cell: (child: DesignerField | DesignerEdge) => (
            <EntityFieldType
                icon={getIconForEntityChild(child)}
                label={child.caption}
                end={<EntityTypeDesignerFieldFlags child={child} />}
            ></EntityFieldType>
        ),
    },
    {
        id: "caption",
        header: "Caption",
        cell: (child: DesignerField | DesignerEdge) => getChildTypeForEntityChild(child),
    },
];

interface ChangesGridProps {
    entityName: string;
    onClose: () => void;
}
const ChangesGrid = (props: ChangesGridProps) => {

    const fullEntityChanges = useEntityDesignerForEntity(props.entityName);

    // Row Actions
    const todoActions = useMemo(
        () => [
            {
                id: "revert",
                render: (field: DesignerField | DesignerEdge) => (
                    <MenuItem
                        id="remove-menu-item"
                        key={`remove-${field.caption}`}
                        onClick={(e) => fullEntityChanges.revertChildChange(field)}
                    >
                        <ListItemIcon>
                            <HistoryIcon />
                        </ListItemIcon>
                        Revert Changes
                    </MenuItem>
                ),
            },
        ],
        [fullEntityChanges.revertChildChange],
    );

    const gridPlugins = useGridPlugins(
        CustomBodyCellContentRenderPlugin,
        usePluginWithParams(PaddingPluggin, {}),
        ColumnsFillRowSpacePlugin,
        usePluginWithParams(RowActionsPlugin, {
            actions: todoActions
        }),
        PinnedColumnsPlugin,
        HideHeaderPlugin,
        usePluginWithParams(EmptyDataPlugin, {
            content: <EmptyMessage />
        }),
    );

    return (
        <Box py={2}>
            <Stack direction="row" justifyContent="space-between" alignItems="center">
                <Stack direction="row" alignItems="center" >
                    <Box p="10px">
                        <EntityChangesIcon designerStatus={fullEntityChanges.nameCaptionEntity?.designerStatus ?? 'unchanged'} />
                    </Box>

                    <Avvvatars value={fullEntityChanges.nameCaptionEntity?.caption ?? ''} style="shape" size={35} />

                    <Typography variant="h5" sx={{ paddingLeft: "12px" }}>
                        {fullEntityChanges.nameCaptionEntity?.caption ?? props.entityName}
                    </Typography>

                    <Typography color="text.secondary" sx={{ paddingLeft: "12px" }}>
                        ( Display Field: {fullEntityChanges.fullDesignerEntity?.displayFieldCaption} )
                    </Typography>
                </Stack>

                <MenuButton
                       slots={{
                        button: (<IconButton><MoreVertIcon /></IconButton>)
                    }}
                    // variant="text"
                    // color="secondary"
                    // actions={[
                    //     {
                    //         type: 'action',
                    //         label: "Revert Changes",
                    //         isCritic: true,
                    //         handler: () => {
                    //             fullEntityChanges.revertEntityChange();
                    //             props.onClose();
                    //         },
                    //         icon: (
                    //             <SvgIcon>
                    //                 <HistoryIcon />
                    //             </SvgIcon>
                    //         ),
                    //     },
                    // ]}
                >

                    <MenuItem data-autoclose onClick={() => { fullEntityChanges.revertEntityChange() }}>
                        <ListItemIcon>
                            <HistoryIcon />
                        </ListItemIcon>
                        <ListItemText>Revert Changes</ListItemText>
                    </MenuItem>
                </MenuButton>
            </Stack>

            <Divider />

            <MosaicDataTable
                plugins={gridPlugins}
                caption="Entity changes"
                items={fullEntityChanges.justChanges?.children ?? []}
                headCells={headerCells}
            />
        </Box>
    );
};
