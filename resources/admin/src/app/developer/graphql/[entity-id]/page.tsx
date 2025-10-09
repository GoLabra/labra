"use client";

import { Button, Container, Link, Stack, SvgIcon, Typography } from "@mui/material";
import { useParams } from "next/navigation";
import { useDocumentTitle } from "@/hooks/use-document-title";
import { useEntities, useFullEntity } from "@/hooks/use-entities";
import Avvvatars from 'avvvatars-react'
import { PageHeader } from "@/shared/components/page-header";
import { useCurrentEntityNameContext } from "@/hooks/use-current-entity";
import { EntityInfo } from "@/features/developer/entity-Info";
import { WithEntityData } from "@/features/developer/with-entity-data";
import { WithEntityDataNewMutation } from "@/features/developer/with-entity-data-new-mutation";
import { WithEntityDataUpdateMutation } from "@/features/developer/with-entity-data-update-mutation";
import { WithEntityDataDeleteMutation } from "@/features/developer/with-entity-data-delete-mutation";
import { ShowGraphQlQuery } from "@/features/developer/show-graph-ql-query";
import { WithEntitySchema } from "@/features/developer/with-entity-schema";

export default function DeveloperEntity() {

    useDocumentTitle({ title: 'Developer' });
    const entityName = useCurrentEntityNameContext()!;
    const fullEntity = useFullEntity({ entityName });

    if (!fullEntity) {
        return (<></>);
    }

    const EntitySchema = WithEntitySchema(ShowGraphQlQuery);
    const EntityData = WithEntityData(ShowGraphQlQuery);
    const EntityDataNewMutation = WithEntityDataNewMutation(ShowGraphQlQuery);
    const EntityDataUpdateMutation = WithEntityDataUpdateMutation(ShowGraphQlQuery);
    const EntityDataDeleteMutation = WithEntityDataDeleteMutation(ShowGraphQlQuery);

    return (
        <>
            <Container
                maxWidth="md"
                sx={{
                    height: '100%',
                    py: 2
                }}>

                <Stack
                    spacing={2}
                    sx={{ height: '100%' }}>

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

                                <Avvvatars value={fullEntity.caption!} style="shape" size={35} />
                                <Typography variant="h1">
                                    {fullEntity.caption}
                                </Typography>

                            </Stack>

                            <Stack
                                direction="row"
                                spacing={1}>

                                {/* <Button
                                size="medium"
                                variant="contained"
                                startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                                onClick={() => addNewEntry()}
                                aria-label="Add new entry"
                                aria-haspopup="dialog">
                                Add New
                            </Button> */}
                            </Stack>

                        </Stack>
                    </PageHeader>

                    <Stack
                        spacing={2}>

                        {/* <EntityInfo entityName={entityName} /> */}

                        <EntitySchema entityName={entityName} />
                        <EntityData entityName={entityName} />
                        <EntityDataNewMutation entityName={entityName} />
                        <EntityDataUpdateMutation entityName={entityName} />
                        <EntityDataDeleteMutation entityName={entityName} />

                    </Stack>

                </Stack>
            </Container>
        </>
    )
}
