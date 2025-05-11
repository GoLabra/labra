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
// import { alpha, darken, lighten } from '@mui/material/styles';
// import type { Components } from '@mui/material/styles/components';

// interface Config {
//     palette: PaletteOptions;
// }

// export const createComponents = ({ palette }: Config): Components => {
//     return {
        
//         MuiAppBar: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.background!.paper
//                 }
//             }
//         },
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
//                     borderColor: palette.neutral![600]
//                 }
//             }
//         },
//         MuiAvatar: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: palette.neutral![800],
//                     color: palette.text!.secondary
//                 }
//             }
//         },
//         MuiButton: {
//             styleOverrides: {
//                 root: {
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
//                         boxShadow: `0px 0px 1px ${palette.neutral![800]}, 0px 1px 3px ${alpha(palette.neutral![900], 0.12)
//                             }`
//                     }
//                 }
//             }
//         },
//         MuiCardHeader: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: darken(palette.background!.paper!, 0.0),
//                 }
//             }
//         },
//         MuiCardFooter: {
//             styleOverrides: {
//                 root: {
//                     backgroundColor: darken(palette.background!.paper!, 0.1),
//                 }
//             }
//         },
//         MuiChip: {
//             styleOverrides: {
//                 avatar: {
//                     backgroundColor: palette.neutral![800]
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
//                     borderColor: palette.neutral![600],
//                     boxShadow: `0px 1px 2px 0px ${alpha(palette.neutral![900], 0.08)}`,
//                     '&:hover': {
//                         backgroundColor: palette.action!.hover
//                     },
//                     [`&.${filledInputClasses.disabled}`]: {
//                         backgroundColor: palette.action!.disabledBackground,
//                         borderColor: palette.neutral![800],
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
//                             fill={palette.background?.paper}
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
//                     backgroundColor: alpha(palette.text!.primary!, 0.1)
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
//                             backgroundColor: alpha(palette.text!.primary!, 0.08)
//                         },
//                         [`& .${switchClasses.thumb}`]: {
//                             backgroundColor: alpha(palette.text!.secondary!, 0.86)
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
//                     //             backgroundColor: lighten(palette.background!.paper!, 0.04)
//                     //         },
//                     //         '.MuiTableCell-actionCell':{
//                     //             backgroundColor: lighten(palette.background!.paper!, 0.04)
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
//                     //     backgroundColor: `${alpha((palette.primary as PaletteColor).main, 0.1)}`,
//                     //     borderRadius: '3px 3px 0 0'
//                     // },
//                     // 'tbody .MuiTableCell-root .column-highlight': {
//                     //     backgroundColor: `${alpha((palette.primary as PaletteColor).main, 0.1)}`,
//                     // }
//                 }
//             }
//         },
//         // MuiTableCell: {
//         //     styleOverrides: {
//         //         root: {
//         //             borderBottomColor: lighten(palette.divider!, 0.7),
//         //         }
//         //     }
//         // },
//         // MuiTableHead: {
//         //     styleOverrides: {
//         //         root: {
//         //             backgroundColor: palette.background!.paper,
//         //             borderBottomColor: lighten(palette.divider!, 0.7),
//         //             [`.${tableCellClasses.root}`]: {
//         //                 color: palette.text!.primary
//         //             }
//         //         }
//         //     }
//         // },
//         // MuiTableRow: {
//         //     styleOverrides: {
//         //         root: {
//         //             [`&.${tableRowClasses.hover}`]: {
//         //                 '&:hover': {
//         //                     backgroundColor: lighten(palette.background!.paper!, 0.04),

//         //                     'td':{
//         //                         backgroundColor: lighten(palette.background!.paper!, 0.04),
//         //                     }
//         //                 }
//         //             }
//         //         }
//         //     }
//         // },
//         MuiToggleButton: {
//             styleOverrides: {
//                 root: {
//                     borderColor: palette.neutral![700]
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
//         },
//         MuiLink: {
//             styleOverrides: {
//                 root: {
//                     color: (palette.primary as PaletteColor).light,
//                     }
//                 }
//             }
//     };
// };
