import { alpha, Box, ButtonBase, PaletteColor, Stack, styled } from "@mui/material";
import EditIcon from "@/assets/icons/iconly/bulk/edit";
import CategoryIcon from "@/assets/icons/iconly/bulk/category";
import { usePathname } from "next/navigation";
import NextLink from 'next/link';
import { paths } from "@/lib/paths";
import { HiExternalLink } from "react-icons/hi";

const MenuItemRoot = styled(Box)({
	padding: '5px 0',
	height: '100%',
	borderWidth: '3px 0',
	borderStyle: 'solid',
	borderColor: 'transparent',

	'&.active': {
		borderBottomColor: 'var(--mui-palette-secondary-main)',
	}
});

interface MenuItemProps {

    icon?: React.ReactNode;
    text: string;
    path: string;
    pathname: string;
    externalLink?: boolean;
}
export const MenuItem = (props: MenuItemProps) => {
    const { icon, text, path, pathname } = props;

    const checkPath = !!(path && pathname);
    const partialMatch = checkPath ? pathname.startsWith(path!) : false;
    const exactMatch = checkPath ? pathname === path : false;

    const active = partialMatch;


    const linkProps = props.externalLink
        ? {
            component: 'a',
            href: path,
            target: '_blank'
        }
        : {
            component: NextLink,
            href: path
        }

    return (
		<MenuItemRoot 
			className={active ? 'active' : ''}>

        <ButtonBase
            {...linkProps}
            sx={{
                alignItems: 'center',
                color: 'text.primary',
                borderRadius: '6px',
                fontFamily: (theme) => theme.typography.fontFamily,
                fontSize: 14,
                fontWeight: 400,
                p: '5px',
				height: '100%',
                whiteSpace: 'nowrap',
				selfAlign: 'middle',
				'&:hover': {
					backgroundColor: (theme) => theme.palette.action.hover
				},
                '&:focus-visible': {
                    backgroundColor: (theme) => theme.palette.action.selected
                },
            }}>

            <Stack direction="row" alignItems="center" gap={1}>
                <Box component="span"
                    sx={{
                        alignItems: 'center',
						opacity: .5,
                        display: 'inline-flex',
                        flexGrow: 0,
                        flexShrink: 0,
                        height: 18,
                        justifyContent: 'center',
                        width: 18,
                        '>svg ': {
                            fontSize: '18px'
                        }
                    }}>
                    {icon}
                </Box>

                <Stack direction="row" alignItems="center" gap="5px">
                    {text}
                    {props.externalLink && <HiExternalLink />}
                </Stack>
            </Stack>

        </ButtonBase>
		</MenuItemRoot>
    )
}

interface TopnavProps {
    direction?: 'row' | 'row-reverse' | 'column' | 'column-reverse';
}
export const Topnav = (props: TopnavProps) => {

    const { direction } = props;
    const pathname = usePathname();

    return (
        <Box component="nav" alignSelf="start">

            <Stack component="ul" margin={0} direction={direction ?? "row"} gap={1} px={0}>
                <Box component="li" sx={{ listStyle: 'none' }}>
                    <MenuItem
                        icon={<EditIcon />}
                        text="Content Manager"
                        path={paths.contentManager.index}
                        pathname={pathname} />
                </Box>
                <Box component="li" sx={{ listStyle: 'none' }}>
                    <MenuItem
                        icon={<CategoryIcon />}
                        text="Entity Type Designer"
                        path={paths.entityTypeDesigner.index}
                        pathname={pathname} />
                </Box>
            </Stack>
        </Box>
    )
}