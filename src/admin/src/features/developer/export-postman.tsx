import { Box, Button, Card, CardContent, CardHeader, Chip as MuiChip, Divider, Stack, styled, Chip, Typography } from '@mui/material';
import { useCallback, useMemo } from 'react';
import { useFullEntities } from '@/hooks/use-entities';
import { GRAPHQL_ENTITY_API_URL, GRAPHQL_QUERY_API_URL } from '@/config/CONST';
import { PostmanCollectionBuilder, PostmanDirectory, PostmanItem } from '@/lib/postman/postman-collection-builder';
import { getEntityDataQuery } from './use-entity-data';
import { getEntityDataDeleteMutationQuery } from './use-entity-data-delete-mutation';
import { getEntityDataNewMutationQuery } from './use-entity-data-new-mutation';
import { getEntityDataUpdateMutationQuery } from './use-entity-data-update-mutation';
import dynamic from 'next/dynamic'
import { entitySchemaQuery } from './use-entity-schema';
import { entitiesSchemaQuery } from './use-entities-schema';
import { saveJsonAs } from '@/lib/utils/open-file';

const saveAs = typeof window === 'undefined' ? undefined : require('save-as').default

export const ExportPostman = () => {

    const fullEntities = useFullEntities({ lazy: false });

    const exportPostmanCollection = useCallback(() => {

        const entities = Object.values(fullEntities.entities);

        const collection = new PostmanCollectionBuilder()
            .addName('labraGO Admin GraphQL')
            .addItem(new PostmanItem()
                .addName('Entities Schema')
                .addUrl(GRAPHQL_ENTITY_API_URL!)
                .addQueryVariables({
                    query: entitiesSchemaQuery,
                    variables: {}
                })
            )
            .addItems(
                entities.map(entity => new PostmanDirectory()
                    .addName(entity.caption!)
                    .addItem(new PostmanItem()
                        .addName(`${entity.caption} Schema`)
                        .addUrl(GRAPHQL_ENTITY_API_URL!)
                        .addQueryVariables({
                            query: entitySchemaQuery,
                            variables: {
                                name: entity.name
                            }
                        })
                    )
                    .addItem(new PostmanItem()
                        .addName('GET')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataQuery(entity))
                    )
                    .addItem(new PostmanItem()
                        .addName('CREATE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataNewMutationQuery(entity.name))
                    )
                    .addItem(new PostmanItem()
                        .addName('UPDATE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataUpdateMutationQuery(entity.name))
                    )
                    .addItem(new PostmanItem()
                        .addName('DELETE')
                        .addUrl(GRAPHQL_QUERY_API_URL!)
                        .addQueryVariables(getEntityDataDeleteMutationQuery(entity.name))
                    )
                )
            )
            .build();

        saveJsonAs('gov-admin-api.postman_collection.json', collection);

    }, [fullEntities]);

    return (<>

        <Card>
            <CardHeader title={(<Stack direction="row" width={1} justifyContent="space-between" alignItems="center">
                <Typography variant="body1"  >Postman</Typography>
                <Button
                    size="medium"
                    variant="outlined"
                    // startIcon={<SvgIcon fontSize="small"><PlusIcon /></SvgIcon>}
                    onClick={() => exportPostmanCollection()}
                    aria-label="Add new entry"
                    aria-haspopup="dialog">
                    Export as Postman Collection
                </Button>
            </Stack>)} />
            <Divider />
            <CardContent>

                <Typography>
                    Postman is a popular API platform that simplifies testing, development, and documentation of APIs. It allows users to send requests, inspect responses, and automate workflows, making it a powerful tool for working with GraphQL and REST APIs.
                </Typography>
            </CardContent>
        </Card>
    </>)
}