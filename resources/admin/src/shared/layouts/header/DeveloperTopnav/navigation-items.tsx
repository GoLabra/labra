import { SvgIcon } from "@mui/material";
import { ReactNode } from "react";
import EditIcon from "@/assets/icons/iconly/bulk/edit";
import CategoryIcon from "@/assets/icons/iconly/bulk/category";

export interface NavBarItem {
    //external?: boolean;
    icon?: ReactNode;
    //items?: NavBarItem[];
    path: string;
    title: string;
}

// export const MainMenuItems: Record<string, NavBarItem> = {
//     ContentManager: {
//         title: 'Content Manager',
//         path: "/content-manager",
//         icon: (
//             <SvgIcon>
//                 <EditIcon />
//             </SvgIcon>
//         )
//     },
//     EntityTypeDesigner: {
//         title: 'Entity Type Designer',
//         path: "/entity-type-designer",
//         icon: (
//             <SvgIcon>
//                 <CategoryIcon />
//             </SvgIcon>
//         )
//     }
// }