import { alpha, Box, ButtonBase, PaletteColor, Stack, styled } from "@mui/material";
import SettingsIcon from "@/assets/icons/iconly/bulk/settings";
import { usePathname } from "next/navigation";
import { paths } from "@/lib/paths";
import { MenuItem } from "../Topnav/top-nav";



interface TopnavProps {
    direction?: 'row' | 'row-reverse' | 'column' | 'column-reverse';
}
export const DevelopTopnav = (props: TopnavProps) => {

    const { direction } = props;
    const pathname = usePathname();


    return (
        <Box component="nav">

            <Stack component="ul" margin={0} direction={direction ?? "row"} gap={1} px={2}>
                <Box component="li" sx={{ listStyle: 'none' }}>
                    <MenuItem
                        icon={<SettingsIcon />}
                        text="Developer"
                        path={paths.developer.index}
                        pathname={pathname} />
                </Box>
                
                <Box component="li" sx={{ listStyle: 'none' }}>
                    <MenuItem
                        icon={<SettingsIcon />}
                        text="Admin"
                        path={paths.contentManager.index}
                        pathname={pathname} />
                </Box>
            </Stack>
        </Box>
    )
}