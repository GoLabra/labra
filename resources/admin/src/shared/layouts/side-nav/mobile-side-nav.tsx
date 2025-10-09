import { Scrollbar } from "@/shared/components/scrollbar";
import { Drawer } from "@mui/material";
import { useState } from "react";

const MOBILE_NAV_WIDTH: number = 280;

interface MobileSideNavProps {
    onClose?: () => void;
    open?: boolean;
    children?: React.ReactNode;
}

export default function MobileSideNavProps(props: MobileSideNavProps) {
    const { open, onClose, children } = props;
    // const [hovered, setHovered] = useState<boolean>(false);

    return (
        <Drawer
            anchor="left"
            open={open}
            onClose={onClose} 

            PaperProps={{
                sx: {
                  width: MOBILE_NAV_WIDTH
                }
              }}>

            <Scrollbar
                sx={{
                    height: '100%',
                    overflowX: 'hidden',
                    '& .simplebar-content': {
                        height: '100%'
                    }
                }}
                tabIndex={-1}>
                {children}
            </Scrollbar>
        </Drawer>
    )
}