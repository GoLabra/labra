import { Box, Chip, Pagination, PaginationItem, Tooltip } from "@mui/material"
import { alpha, styled } from '@mui/material/styles';

export const CenterPagination = styled(Box)({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: '10px'
});

interface CoolPaginationProps {
    page: number;
    pagesCount: number;
    totalItems: number;
    onChange: (page: number) => void
}
export const CoolPagination = (props: CoolPaginationProps) => {
    const { page, pagesCount, totalItems, onChange } = props;

    return (<Pagination
        count={pagesCount}
        page={page}
        onChange={(_, value: number) => onChange(value)}
        renderItem={(item) => (
            <>
                {(item.type === 'next' && totalItems > 0) && (
                    <Tooltip title="Total Items Count">
                        <Chip label={totalItems} 
                        
                        sx={{
                            minWidth: '40px',
                            padding: '0 10px',
                            backgroundColor: (theme) => theme.palette.mode === 'dark' ? alpha(theme.palette.neutral[900], .6) : alpha(theme.palette.neutral[50], .2),
                            margin: '0 10px',
                            borderRadius: 1
                        }} />
                    </Tooltip>
                )}
                <PaginationItem {...item} />
            </>
        )} 
    />)
}