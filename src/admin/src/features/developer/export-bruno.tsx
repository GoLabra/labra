import { Box, Button, Card, CardContent, CardHeader, Chip as MuiChip, Divider, Stack, styled, Chip, Typography } from '@mui/material';
import { useCallback, useMemo } from 'react';
import { useFullEntities } from '@/hooks/use-entities';
import { GRAPHQL_ENTITY_API_URL, GRAPHQL_QUERY_API_URL } from '@/config/CONST';
import { PostmanCollectionBuilder, PostmanDirectory, PostmanItem } from '@/lib/postman/postman-collection-builder';
import { getEntityDataQuery } from './use-entity-data';
import { getEntityDataDeleteMutationQuery } from './use-entity-data-delete-mutation';
import { getEntityDataNewMutationQuery } from './use-entity-data-new-mutation';
import { getEntityDataUpdateMutationQuery } from './use-entity-data-update-mutation';
import { BrunoCollectionBuilder, BrunoDirectory, BrunoItem } from '@/lib/bruno/bruno-collection-builder';
import dynamic from 'next/dynamic'
import { entitiesSchemaQuery } from './use-entities-schema';
import { entitySchemaQuery } from './use-entity-schema';
import { saveJsonAs } from '@/lib/utils/open-file';

export const ExportBruno = () => {

    const fullEntities = useFullEntities({ lazy: false });

    const exportPostmanCollection = useCallback(() => {

        const entities = Object.values(fullEntities.entities);

        const collection = new BrunoCollectionBuilder()
            .addName('labraGO Admin GraphQL')
            .addItem(new BrunoItem()
                .addName('Entities Schema')
                .addUrl(GRAPHQL_ENTITY_API_URL!)
                .addQueryVariables({
                    query: entitiesSchemaQuery,
                    variables: {}
                })
            )
            .addItems(
                entities.map(entity => new BrunoDirectory()
                    .addName(entity.caption!)
                    .addItem(new BrunoItem()
                        .addName(`${entity.caption} Schema`)
                        .addUrl(GRAPHQL_ENTITY_API_URL!)
                        .addQueryVariables({
                            query: entitySchemaQuery,
                            variables: {
                                name: entity.name
                            }
                        })
                    )
                    .addItem(new BrunoItem()
                        .addName('GET')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataQuery(entity))
                    )
                    .addItem(new BrunoItem()
                        .addName('CREATE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataNewMutationQuery(entity.name))
                    )
                    .addItem(new BrunoItem()
                        .addName('UPDATE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataUpdateMutationQuery(entity.name))
                    )
                    .addItem(new BrunoItem()
                        .addName('DELETE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataDeleteMutationQuery(entity.name))
                    )
                )
            )
            .build();

        saveJsonAs('gov-admin-api.bruno_collection.json', collection);

    }, [fullEntities]);

    return (<>

        <Card>
            <CardHeader title={(<Stack direction="row" width={1} justifyContent="space-between" alignItems="center">
                <Typography variant="body1">Bruno</Typography>
                <Button
                    size="medium"
                    variant="outlined"
                    // startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                    onClick={() => exportPostmanCollection()}
                    aria-label="Add new entry"
                    aria-haspopup="dialog">
                    Export as Bruno Collection
                </Button>
            </Stack>)} />
            <Divider />
            <CardContent>
                <Typography>
                    Bruno is a popular API platform that simplifies testing, development, and documentation of APIs. It allows users to send requests, inspect responses, and automate workflows, making it a powerful tool for working with GraphQL and REST APIs.
                </Typography>
            </CardContent>

        </Card>
    </>)
}