

import { Box, Divider, Drawer, Stack } from "@mui/material";
import { usePathname } from "next/navigation";
import { useState } from "react";
import { SideNavItem } from "./side-nav-item";
import { NavBarItem } from "./../header/Topnav/navigation-items";
import { Header } from "../header/Header";
import { Scrollbar } from "@/shared/components/scrollbar";

const SIDE_NAV_WIDTH: number = 240;
const TOP_NAV_HEIGHT: number = 64;

interface SideNavProps {
    pinned?: boolean;
    onPin?: () => void;
    children?: React.ReactNode;
}

export default function SideNav(props: SideNavProps) {
    const { onPin, pinned = false, children } = props;
    // const [hovered, setHovered] = useState<boolean>(false);
    // //const pathname = usePathname();
    // const collapse = !(pinned || hovered);

    return (
        //     <Drawer
        //   open
        //   variant="permanent"
        //   PaperProps={{
        //     onMouseEnter: () => { setHovered(true); },
        //     onMouseLeave: () => { setHovered(false); },
        //     sx: {
        //       backgroundColor: 'background.default',
        //       height: `calc(100% - ${TOP_NAV_HEIGHT}px)`,
        //       overflowX: 'hidden',
        //       top: TOP_NAV_HEIGHT,
        //       transition: 'width 250ms ease-in-out',
        //       width: SIDE_NAV_WIDTH,
        //       zIndex: (theme) => theme.zIndex.appBar - 100,
        //       boxShadow: 'none'
        //     }
        //   }}
        // >
        <Box component="aside"  sx={{
            top: TOP_NAV_HEIGHT,
            width: SIDE_NAV_WIDTH,
            gridArea: 'nav',
            borderRight: '1px solid var(--mui-palette-divider)',
        }}>
            <Box sx={{
                overflowX: 'hidden',
                height: `calc(100vh - ${TOP_NAV_HEIGHT}px)`,
                position: 'sticky',
                top: TOP_NAV_HEIGHT,
            }}>


            <Scrollbar
                sx={{
                    width: '100%',
                    height: '100%',
                    overflowX: 'hidden',
                    '& .simplebar-content': {
                        height: '100%'
                    }
                }}
                tabIndex={-1}>

                {children}
            </Scrollbar>


            </Box>

            {/* <Divider /> */}
            {/* <Header /> */}
        </Box>
    )
}