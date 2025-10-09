// import styled from "@emotion/styled";

// export const DialogStackVisual = ({ count }: { count: number }) => {

//     if (count == 0) {
//         return null;
//     }

//     const width = 5 + (1 / count) * 10;

//     //TODO: get the colors from the theme
//     const StackElement = styled('div')(({ theme }) => ({
//         content: "''",
//         width: width,
//         backgroundColor: 'rgb(74 85 113)',
//         position: 'absolute',
//         left: -width,
//         top: 10,
//         bottom: 32,
//         height: '95%',
//         borderRadius: 5,
//         borderTopRightRadius: 0,
//         borderBottomRightRadius: 0,
//         background: 'linear-gradient(90deg, rgb(41 48 66) 0%, rgb(31 36 50) 100%)',
//     }));

//     return <>
//         <StackElement>
//             <DialogStackVisual count={count - 1}>
//             </DialogStackVisual>
//         </StackElement>
//     </>
// }