// with-entity-data.tsx

import { ComponentType, FC } from "react";
import { ShowGraphQlQueryProps } from "./show-graph-ql-query";
import { useFullEntity } from "@/hooks/use-entities";
import { GqlDataQueryBuilder } from "@/lib/apollo/builders/gqlQueryBuilder";
import { Box, Card, CardContent, Chip, Divider, Stack } from "@mui/material";

interface EntityInfoProps {
    entityName: string;
}
export const EntityInfo = (props: EntityInfoProps) => {

    const fullEntity = useFullEntity({ entityName: props.entityName });

    return (<>

        <Card>

            <Divider />
            <CardContent>
                == FIELDS ====<br />
                <Stack direction="row" gap={1}>
                    {fullEntity?.fields.map(field => (
                        <Chip key={field.name} component="span" label={field.name} color="secondary" variant="outlined" size="small" />
                    ))}
                </Stack>
            </CardContent>
            <Divider />
            <CardContent>
                == EDGES ====<br />
                <Stack direction="row" gap={1}>
                    {fullEntity?.edges.map(edge => (
                        <Chip key={edge.name} component="span" label={edge.name} color="secondary" variant="outlined" size="small" />
                    ))}
                </Stack>
            </CardContent>
        </Card>

    </>)
}