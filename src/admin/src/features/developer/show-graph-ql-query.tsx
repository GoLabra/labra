import { Box, Button, Card, CardContent, CardHeader, Chip as MuiChip, Divider, Stack, styled, IconButton } from '@mui/material';
import gqlPrettier from 'graphql-prettier';
import copy from 'copy-to-clipboard';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import { useCallback, useMemo, useState } from 'react';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import { gql, useApolloClient } from '@apollo/client';
import { GRAPHQL_ENTITY_PLAYGROUND_URL, GRAPHQL_QUERY_PLAYGROUND_URL } from '@/config/CONST';
import { ENTITY_CONTEXT } from '@/lib/apollo/apolloWrapper';

const Chip = styled(MuiChip)(({ theme }) => ({
    borderRadius: '3px',
}));
export interface ShowGraphQlQueryProps {
    title: string
    query: string;
    variables: Record<string, any> | null;
    context: any | null;
}
export const ShowGraphQlQuery: React.FC<ShowGraphQlQueryProps> = (props: ShowGraphQlQueryProps) => {

    const client = useApolloClient();

    const query = useMemo(() => gqlPrettier(props.query), [props.query]);
    const variables = useMemo(() => props.variables ? JSON.stringify(props.variables, null, 2) : null, [props.variables]);

    const [result, setResult] = useState<any>(null);

    const runQuery = useCallback(() => {
        client.query({ query: gql(query), variables: props.variables ?? {}, fetchPolicy: "network-only", context: props.context }).then((response) => {
            setResult(response);
        }).catch((error) => {
            setResult(error);
        });

    }, [client, query, variables, props.context]);

    const runQueryInPlayground = useCallback(() => {
        const playgroundUrl = props.context ? GRAPHQL_ENTITY_PLAYGROUND_URL : GRAPHQL_QUERY_PLAYGROUND_URL;
        window.open(`${playgroundUrl!}?q=${encodeURIComponent(query)}&v=${encodeURIComponent(variables ?? '{}')}`, '_blank');
    }, [query, variables, props.context]);

    return (<>

        <Card>
            <CardHeader title={props.title} />
            <Divider />
            <CardContent>

                <Stack direction="row" width={1} justifyContent="space-between">
                    <Chip label="QUERY" color="success" variant="outlined" size="small" />
                    <Button variant="outlined" onClick={() => copy(query)}>
                        <ContentCopyIcon fontSize='small' />
                    </Button>
                </Stack>

                <Box component="pre" sx={{ whiteSpace: 'pre-wrap', maxHeight: '200px', overflow: 'auto' }}>
                    {query}
                </Box>
            </CardContent>
            <Divider />
            
            {props.variables && (<CardContent>
                <Stack direction="row" width={1} justifyContent="space-between">
                    <Chip label="VARIABLES" color="success" variant="outlined" size="small" />
                    <Button variant="outlined" onClick={() => copy(variables!)}>
                        <ContentCopyIcon fontSize='small' />
                    </Button>
                </Stack>

                <Box component="pre" sx={{ whiteSpace: 'pre-wrap', maxHeight: '200px', overflow: 'auto' }}>
                    {variables}
                </Box>
            </CardContent>)}

            <Divider />
            <CardContent>
                <Stack direction="row" width={1} justifyContent="space-between"> 
                    <Box>
                        {result && (<Chip label="RESULT" color="success" variant="outlined" size="small" />)}
                    </Box>

                    <Stack direction="row" gap={1}>

                        {result && (<Button variant="outlined" onClick={() => copy(JSON.stringify(result, null, 2))}>
                            <ContentCopyIcon fontSize='small' />
                        </Button>)}

                        <Button variant="outlined" onClick={() => runQuery()}>
                            <PlayArrowIcon />
                        </Button>

                        <Button variant="outlined" onClick={() => runQueryInPlayground()}>
                            Run in Playground
                        </Button>
                    </Stack>
                </Stack>

                {result && (<Box component="pre" sx={{ whiteSpace: 'pre-wrap', maxHeight: '200px', overflow: 'auto' }}>
                    {JSON.stringify(result, null, 2)}
                </Box>)}

            </CardContent>
        </Card>
    </>)
}
ShowGraphQlQuery.displayName = 'ShowGraphQlQuery';