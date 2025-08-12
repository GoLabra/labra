import { Box, Button, ButtonBase, Chip, Stack, Typography } from "@mui/material";
import { ColumnDef, Filter, FilterValue } from "mosaic-data-table";
import { ReactNode, useCallback, useMemo } from "react";
import CancelIcon from '@mui/icons-material/Cancel';
import DeleteIcon from '@mui/icons-material/Delete';

interface FilterEditorProps {
    filter: Filter
    onRemove: (key?: string) => void;

    headCells?: ColumnDef[];
}
export const FilterEditor = (props: FilterEditorProps) => {

    const entries = useMemo(() => Object.entries(props.filter), [props.filter]);

    const headCellsMap = useMemo(() => {
        return props.headCells?.reduce((acc: Record<string, ColumnDef>, headCell: ColumnDef) => {
            acc[headCell.id] = headCell;
            return acc;
        }, {});
    }, [props.headCells]);

    const onClear = useCallback(() => {
        props.onRemove();
    }, [props.onRemove]);

    const onDelete = useCallback((key: string) => {
        props.onRemove(key);
    }, [props.onRemove]);

	const isVisible = !!Object.keys(props.filter).length;
	if(!isVisible){
		return null;
	}


    return (
        <Stack direction="row" gap={2} alignItems="center" flexWrap={'wrap'} sx={{
			marginTop: 2
		}}>

            {entries.map(([key, value]) => (
                <FilterEditorItem key={key} name={key} filter={value} onDelete={() => onDelete(key)} headCell={headCellsMap?.[key]} />
            ))}

            {entries.length > 0 &&
                <Button startIcon={<DeleteIcon />} onClick={onClear} >
                    Clear
                </Button>}
        </Stack>
    );
}

interface FilterEditorItemProps {
    name: keyof Filter;
    filter: FilterValue;
    onDelete?: () => void;

    headCell?: ColumnDef;
}
export const FilterEditorItem = (props: FilterEditorItemProps) => {

    const name = useMemo(() => {
        return typeof props.headCell?.header === 'function'
                    ? props.headCell?.header()
                    : props.headCell?.header ?? props.name;
    }, [props.headCell, props.name]);

    const value = useMemo(() => {
        return props.headCell?.cell?.({[props.name]:props.filter.value}) ?? props.filter.value; return
    }, [props.headCell, props.filter]);

    return (
        <Box sx={{
            borderWidth: '1.5px',
            borderStyle: 'dashed',
            borderColor: 'divider',
            padding: '3px 3px 3px 10px',
            borderRadius: '6px',
        }}>
            <Stack direction="row" gap={1} alignItems="center" flexWrap={'wrap'}>
                <Typography variant="subtitle1">
                    {name}:
                </Typography>

                <Chip label={value} deleteIcon={<CancelIcon />} onDelete={props.onDelete} sx={{
                    '& .MuiChip-deleteIcon': {
                        color: 'text.secondary'  // or any other color like '#FF0000' or 'red'
                    }
                }} />
                {/* <ButtonBase >
                </ButtonBase> */}
            </Stack>
        </Box>
    );
}