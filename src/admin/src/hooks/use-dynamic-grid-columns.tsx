import { localeConfig } from "@/config/locale-config";
import { Edge, Field, RelationType } from '@/lib/apollo/graphql.entities';
import { ApiFieldTypes } from '@/types/field-type-descriptor';
import CheckIcon from '@mui/icons-material/Check';
import { Button, Chip, Stack, Typography } from '@mui/material';
import { ColumnDef } from 'mosaic-data-table';
import { MutableRefObject, ReactNode, useMemo, useRef } from 'react';
// import { stringAvatar } from '@/lib/utils/avatar';
import { stringToDate, stringToDateTime, stringToTime } from '@/core-features/dynamic-form/value-convertor';

export type ColumnOptions = {
    hasSort?: boolean,
    width?: number,
}

export type ColumnSourceMeta =
    | { kind: 'field'; field: Field }
    | { kind: 'edge'; edge: Edge };

interface UseDynamicGridColumnsProps {
    entityName: string,
    fields?: Field[],
    edges?: Edge[],
    displayFieldName?: string,
	openRelation: (entityName: string, edge: Edge, entryId: string) => void,
    showId: boolean,
    transformColumn?: (column: ColumnDef, meta: ColumnSourceMeta) => ColumnDef | void,
}
export const useDynamicGridColumns = ({
    entityName,
    fields = [],
    edges = [],
    displayFieldName,
    openRelation,
    showId,
    transformColumn
}: UseDynamicGridColumnsProps): ColumnDef[] => {

    const expansionStoreRef = useRef(openRelation);
    expansionStoreRef.current = openRelation;

    return useMemo(() => {
        const fieldColumns: ColumnDef[] = fields
            .filter(i => showId ? true : i.name != 'id')
            .filter(i => !i.private) // TODO: this should be hidden from BE
            .filter(i => i.type != 'RichText')
			.filter(i => i.type != 'LongText') 
            .filter(i => i.type != 'Json')
            .map(field => {
                let col = fieldToColumn(field);
                if (transformColumn) {
                    col = transformColumn(col, { kind: 'field', field }) ?? col;
                }
                return col;
            });

        const edgeColumns: ColumnDef[] = edges.map(edge => {
            let col = edgeToColumn(entityName, edge, expansionStoreRef);
            if (transformColumn) {
                col = transformColumn(col, { kind: 'edge', edge }) ?? col;
            }
            return col;
        });

        const builtColumns: ColumnDef[] = [
            ...fieldColumns,
            ...edgeColumns,
        ];

        return builtColumns;
    }, [showId, displayFieldName, fields, edges, transformColumn, entityName]);
}


const fieldToColumn = (field: Field): ColumnDef<Field> => {

    // Add avatar for name column
    switch (field.type as ApiFieldTypes) {
        case 'ShortText': return shortTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'RichText': return richTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Integer': return integerColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'DateTime': return dateTimeColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Date': return dateColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Time': return timeColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'Boolean': return booleanColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
        case 'MultipleChoice': return multiChoiceColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: false });
        default: return shortTextColumnDef(field.name, field.caption, (row: any) => row[field.name], { hasSort: true });
    }
}

const edgeToColumn = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>): ColumnDef<Edge> => {

    switch (edge.relationType) {
        case RelationType.One:
        case RelationType.OneToOne:
        case RelationType.ManyToOne:
            return oneColumnDef(entityName, edge, openRelation);
        case RelationType.Many:
		case RelationType.OneToMany:
        case RelationType.ManyToMany:
            return manyColumnDef(entityName, edge, openRelation);
    }

    return {
        id: edge.name,
        header: '',
        cell: (row: any) => ''
    }
}

const shortTextColumnDef = (name: string, caption: string, render: (row: any) => ReactNode, options?: ColumnOptions): ColumnDef<any> => {

    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const richTextColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 300,
        hasSort: options?.hasSort ?? false
    };
}

const integerColumnDef = (name: string, caption: string, render: (row: any) => number, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => render(row),
        width: options?.width ?? 150,
        hasSort: options?.hasSort ?? false
    };
}

const dateTimeColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

			return (<>
				<Typography variant="body2" component='span' color='var(--mui-palette-text-secondary)'>{stringToDateTime(value)?.format(localeConfig.date.displayFormat) ?? undefined}, </Typography>
				<Typography variant="body2" component='span'>{stringToDateTime(value)?.format(localeConfig.time.displayFormat) ?? undefined}</Typography>
			</>)
            return stringToDateTime(value)?.format(localeConfig.dateTime.displayFormat) ?? undefined;

        },
        width: options?.width ?? 210,
        hasSort: options?.hasSort ?? false
    };
}

const dateColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return stringToDate(value)?.format(localeConfig.date.displayFormat) ?? undefined;

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const timeColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return stringToTime(value)?.format(localeConfig.time.displayFormat) ?? undefined;

        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const booleanColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);
            if (!value) {
                return '';
            }

            return (<CheckIcon />)

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const multiChoiceColumnDef = (name: string, caption: string, render: (row: any) => string, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: name,
        header: caption,
        cell: (row: any) => {

            const value = render(row);

            if (!Array.isArray(value)) {
                return value;
            }

            return (<Stack direction="row" gap="3px">
                {value.map((i) => {
                    return (<Chip key={i} label={i} size="small" variant="outlined" />)
                })}
            </Stack>
            )

        },
        width: options?.width ?? 120,
        hasSort: options?.hasSort ?? false
    };
}

const oneColumnDef = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>, options?: ColumnOptions): ColumnDef<any> => {

    return {
        id: edge.name,
        header: edge.caption,
        cell: (row: any) => {
            const entityValue = row[edge.name];
            if (!entityValue) {
                return '';
            }

            const displayValue = entityValue[edge.relatedEntity.displayField.name] || `Id: ${entityValue['id']}`;

            return <Button
                variant="text"
                sx={{
                    padding: 0
                }}
                onClick={() => {
					openRelation.current(entityName, edge, row.id);
                }}>
                <Typography
                    color="text.primary"
                    variant='body2'
                    sx={{ textDecoration: 'underline' }}>{displayValue}</Typography></Button>
        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

const manyColumnDef = (entityName: string, edge: Edge, openRelation: MutableRefObject<(entityName: string, edge: Edge, entryId: string) => void>, options?: ColumnOptions): ColumnDef<any> => {
    return {
        id: edge.name,
        header: edge.caption,
        cell: (row: any) => {

            return <Button
                variant="text"
                sx={{
                    padding: 0
                }}
                onClick={() => {

					openRelation.current(entityName, edge, row.id);

                }}>
                <Typography
                    color="text.primary"
                    variant='body2'
                    sx={{ textDecoration: 'underline' }}>VIEW</Typography>
            </Button>
        },
        width: options?.width ?? 180,
        hasSort: options?.hasSort ?? false
    };
}

