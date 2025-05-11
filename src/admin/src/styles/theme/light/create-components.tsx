// import type { PaletteColor, PaletteOptions } from '@mui/material';
// import {
//     filledInputClasses,
//     paperClasses,
//     radioClasses,
//     switchClasses,
//     tableCellClasses,
//     tableRowClasses
// } from '@mui/material';
// import { common } from '@mui/material/colors';
// import { alpha, lighten, darken } from '@mui/material/styles';
// import type { Components } from '@mui/material/styles/components';

// interface Config {
//     palette: PaletteOptions;
// }

// export const createComponents = ({ palette }: Config): Components => {
//     return {
//         // MuiAppBar: {
//         //     styleOverrides: {
//         //         root: {
//         //             backgroundColor: palette.background!.paper,
//         //         }
//         //     }
//         // },
//         MuiDrawer: {
//             styleOverrides: {
//                 paper: {
//                     backgroundColor: palette.background?.default
//                 }
//             }
//         },
//         MuiAutocomplete: {
//             styleOverrides: {
//                 paper: {
//                     borderWidth: 1,
//                     borderStyle: 'solid',
//                     borderColor: palette.divider
//                 }
//             }
//         },
//         MuiAvatar: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.neutral![200],
//                     color: palette.text!.secondary
//                 }
//             }
//         },
//         MuiButton: {
//             styleOverrides: {
//                 contained: {
//                     '&:focus': {
//                         boxShadow: `${alpha((palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
//                     }
//                 }
//             }
//         },
//         MuiCard: {
//             styleOverrides: {
//                 root: {
//                     [`&.${paperClasses.elevation1}`]: {
//                         boxShadow: `0px 0px 1px ${palette.neutral![200]}, 0px 1px 3px ${alpha(palette.neutral![800], 0.08)
//                             }`
//                     }
//                 }
//             }
//         },
//         MuiCardHeader: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: darken(palette.background!.paper!, 0.015),
//                 }
//             }
//         },
//         MuiCardFooter: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: darken(palette.background!.paper!, 0.015),
//                 }
//             }
//         },
//         MuiChip: {
//             styleOverrides: {
//                 avatar: {
//                     backgroundColor: palette.neutral![200]
//                 }
//             }
//         },
//         MuiInputBase: {
//             styleOverrides: {
//                 input: {
//                     '&::placeholder': {
//                         color: palette.text!.secondary
//                     }
//                 }
//             }
//         },
//         MuiFilledInput: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.background!.paper,
//                     borderColor: palette.neutral![300],
//                     boxShadow: `0px 1px 2px 0px ${alpha(palette.neutral![800], 0.08)}`,
//                     '&:hover': {
//                         backgroundColor: palette.action!.hover
//                     },
//                     [`&.${filledInputClasses.disabled}`]: {
//                         backgroundColor: palette.action!.disabledBackground,
//                         borderColor: palette.neutral![300],
//                         boxShadow: 'none'
//                     },
//                     [`&.${filledInputClasses.focused}`]: {
//                         backgroundColor: 'transparent',
//                         borderColor: (palette.primary as PaletteColor).main,
//                         boxShadow: `${alpha((palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
//                     },
//                     [`&.${filledInputClasses.error}`]: {
//                         borderColor: (palette.error as PaletteColor).main,
//                         boxShadow: `${alpha((palette.error as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
//                     }
//                 }
//             }
//         },
//         MuiFormLabel: {
//             styleOverrides: {
//                 root: {
//                     color: palette.text!.primary
//                 }
//             }
//         },
//         MuiRadio: {
//             defaultProps: {
//                 checkedIcon: (
//                     <svg
//                         width="18"
//                         height="18"
//                         viewBox="0 0 18 18"
//                         fill="none"
//                         xmlns="http://www.w3.org/2000/svg"
//                     >
//                         <rect
//                             width="18"
//                             height="18"
//                             rx="9"
//                             fill="currentColor"
//                         />
//                         <rect
//                             x="2"
//                             y="2"
//                             width="14"
//                             height="14"
//                             rx="7"
//                             fill="currentColor"
//                         />
//                         <rect
//                             x="5"
//                             y="5"
//                             width="8"
//                             height="8"
//                             rx="4"
//                             fill={palette.background!.paper}
//                         />
//                     </svg>
//                 ),
//                 icon: (
//                     <svg
//                         width="18"
//                         height="18"
//                         viewBox="0 0 18 18"
//                         fill="none"
//                         xmlns="http://www.w3.org/2000/svg"
//                     >
//                         <rect
//                             width="18"
//                             height="18"
//                             rx="9"
//                             fill="currentColor"
//                         />
//                         <rect
//                             x="2"
//                             y="2"
//                             width="14"
//                             height="14"
//                             rx="7"
//                             fill={palette.background!.paper}
//                         />
//                     </svg>
//                 )
//             },
//             styleOverrides: {
//                 root: {
//                     color: palette.text!.secondary,
//                     [`&:hover:not(.${radioClasses.checked})`]: {
//                         color: palette.text!.primary
//                     }
//                 }
//             }
//         },
//         MuiSkeleton: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.neutral![100]
//                 }
//             }
//         },
//         MuiSwitch: {
//             styleOverrides: {
//                 root: {
//                     '&:focus-within': {
//                         boxShadow: `${alpha((palette.primary as PaletteColor).main, 0.25)} 0 0 0 0.2rem`
//                     }
//                 },
//                 switchBase: {
//                     [`&.${switchClasses.disabled}`]: {
//                         [`&+.${switchClasses.track}`]: {
//                             backgroundColor: alpha(palette.neutral![500], 0.08)
//                         },
//                         [`& .${switchClasses.thumb}`]: {
//                             backgroundColor: alpha(palette.neutral![500], 0.26)
//                         }
//                     }
//                 },
//                 track: {
//                     backgroundColor: palette.neutral![500]
//                 },
//                 thumb: {
//                     backgroundColor: common.white
//                 }
//             }
//         },
//         MosaicDataTable: {
//             styleOverrides: {
//                 root: {
//                     // 'tbody tr.MuiTableRow-hover ': {
//                     //     '&:hover': {
//                     //         '.MuiTableCell-selectionCell':{
//                     //             backgroundColor: darken(palette.background!.paper!, 0.04)
//                     //         },
//                     //         '.MuiTableCell-actionCell':{
//                     //             backgroundColor: darken(palette.background!.paper!, 0.04)
//                     //         }
//                     //     }
//                     // },
//                     // 'tbody tr.MuiTableRow-root.highlight ': {
//                     //     backgroundColor: `color-mix(in srgb, ${(palette.primary as PaletteColor).main} 10%, ${palette.background!.paper} 90%)`,

//                     //     '.MuiTableCell-selectionCell':{
//                     //         backgroundColor: `color-mix(in srgb, ${(palette.primary as PaletteColor).main} 10%, ${palette.background!.paper} 90%)`,
//                     //     },
//                     //     '.MuiTableCell-actionCell':{
//                     //         backgroundColor: `color-mix(in srgb, ${(palette.primary as PaletteColor).main} 10%, ${palette.background!.paper} 90%)`,
//                     //     }
//                     // },
//                     // 'thead .MuiTableCell-root .column-highlight': {
//                     //     backgroundColor: `${alpha((palette.primary as PaletteColor).main, 0.2)}`,
//                     //     borderRadius: '3px 3px 0 0'
//                     // },
//                     // 'tbody .MuiTableCell-root .column-highlight': {
//                     //     backgroundColor: `${alpha((palette.primary as PaletteColor).main, 0.2)}`,
//                     // }
//                 }
//             }
//         },
//         // MuiTableHead: {
//         //     styleOverrides: {
//         //         root: {
//         //             backgroundColor: palette.background!.paper,
//         //             borderBottomColor: darken(palette.divider!, 0.3),
//         //             [`.${tableCellClasses.root}`]: {
//         //                 color: palette.text!.secondary
//         //             }
//         //         }
//         //     }
//         // },
//         // MuiTableCell: {
//         //     styleOverrides: {
//         //         root: {
//         //             borderBottomColor: darken(palette.divider!, 0.3),
//         //         }
//         //     }
//         // },
//         // MuiTableRow: {
//         //     styleOverrides: {
//         //         root: {
//         //             [`&.${tableRowClasses.hover}`]: {
//         //                 '&:hover': {
//         //                     backgroundColor: darken(palette.background!.paper!, 0.04)
//         //                 }
//         //             }
//         //         }
//         //     }
//         // },
//         MuiToggleButton: {
//             styleOverrides: {
//                 root: {
//                     borderColor: palette.divider
//                 }
//             }
//         },
//         MuiDialogTitle: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.background?.default
//                 }
//             }
//         },
//         MuiDialogActions: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.background?.default
//                 }
//             }
//         },
//         MuiStepContent: {
//             styleOverrides: {
//                 root: {
//                     borderLeftColor: `${alpha((palette.primary as PaletteColor).main, 0.26)}`,

//                 }
//             }
//         },
//         MuiStepConnector: {
//             styleOverrides: {
//                 root: {

//                     '>.MuiStepConnector-line': {
//                         borderLeftColor: `${alpha((palette.primary as PaletteColor).main, 0.26)}`,
//                     }
//                 }
//             }
//         }
//     };
// };
