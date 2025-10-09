import { Box, Collapse, Container, Divider, Drawer, Paper, styled } from "@mui/material";
import { useState } from "react";
import { usePathname } from "next/navigation";
import { Scrollbar } from "./scrollbar";

interface VerticalScrollPanelProps {
    width?: number;
    children: React.ReactNode,
}

export default function VerticalScrollPanel(props: VerticalScrollPanelProps) {

    const { children } = props;
    const width = props.width ?? 240;

    //const [SIDE_NAV_WIDTH, setSIDE_NAV_WIDTH] = useState<number>(width);
    const SIDE_NAV_WIDTH = 240;

    return (
        <Box
            width={SIDE_NAV_WIDTH}
            minWidth={SIDE_NAV_WIDTH}>

            <Box
                position="fixed"
                width={SIDE_NAV_WIDTH}
                p={0}
                sx={{
                    top: 0,
                    bottom: 0,
                    borderWidth: 0,
                    borderRightWidth: 1,
                    borderStyle: 'solid',
                    borderColor: 'divider',

                }}>



                <Scrollbar
                    sx={{
                        
                        height: '100%',
                        overflowX: 'hidden',
                        '& .simplebar-content': {
                            height: '100%'
                        },
                        maxHeight: `100vh`
                    }}
                    tabIndex={-1}>
                    {/* {Array.from(Array(100)).map((v, i) => (
                        <div key={i}>Testing {i}</div>
                    ))} */}

                    {children}

                </Scrollbar>

            </Box>
        </Box>
    )
}