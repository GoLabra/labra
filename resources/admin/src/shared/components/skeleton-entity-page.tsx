import { styled, alpha, Stack, Typography, Skeleton, Card, CardContent } from "@mui/material";
import { Box } from "@mui/material";
import { PageHeader } from "./page-header";
import Avvvatars from "avvvatars-react";
import { ColumnsFillRowSpacePlugin, MosaicDataTable, PinnedColumnsPlugin, RowActionsPlugin, SkeletonLoadingPlugin, useGridPlugins, usePluginWithParams } from "mosaic-data-table";

export const SkeletonEntityPage = () => {

    const gridPlugins = useGridPlugins(
        ColumnsFillRowSpacePlugin,
        usePluginWithParams(RowActionsPlugin, {
            actions: []
        }),
        usePluginWithParams(SkeletonLoadingPlugin, {
            isLoading: true,
        }),
    );
    return (
        <>
            <PageHeader
                sx={{
                    pl: 1
                }}>
                <Stack
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                    spacing={1}>

                    <Stack direction="row" alignItems="center" gap={1}>
                        <Skeleton variant="circular" width={35} height={35} />
                        <Typography variant="h2">
                            <Skeleton width={200} />
                        </Typography>
                    </Stack>

                    <Stack
                        direction="row"
                        spacing={1}>


                        <Skeleton sx={{ height: 30, width: 100, borderRadius: 1 }} variant="rectangular" />
                        <Skeleton sx={{ height: 30, width: 100, borderRadius: 1 }} variant="rectangular" />

                    </Stack>

                </Stack>

               
            </PageHeader>

            <Card>
                    <CardContent>
                        <MosaicDataTable
                            plugins={gridPlugins}
                            caption="Entity changes"
                            items={[]}
                            headCells={[]}
                        />
                    </CardContent>
                </Card>
        </>
    );
}