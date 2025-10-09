// import { Edge, Field } from "@/lib/apollo/graphql";

// const convertDataToGraphQl = (formData: any, fields?: Field[], edges?: Edge[]) => {

//     // set fields
//     let result = fields?.reduce((acc: any, item: any) => {
//         acc[item.name] = formData[item.name]
//         return acc
//     }, {});

//     // set edges
//     result = edges?.reduce((acc: any, item: any) => {
//         const edgeValue = formData[item.name];
//         if (!edgeValue) {
//             return acc;
//         }

//         switch (edgeValue.valueType) {
//             case "connect":
//                 acc[item.name] = {
//                     connect: {
//                         id: edgeValue.value
//                     }
//                 };
//                 break;

//             case "create":
//                 acc[item.name] = {
//                     create: {
//                         ...edgeValue.value
//                     }
//                 };
//                 break;
//         }

//         return acc;
//     }, result);

//     return result;

// };

// const convertDataFromGraphQl = (graphQlData: any, fields?: Field[], edges?: Edge[]) => {

//     if(!graphQlData){
//         return graphQlData;
//     }

//     // set fields
//     let result = fields?.reduce((acc: any, item: any) => {
//         acc[item.name] = graphQlData[item.name]
//         return acc
//     }, {});

//     // set edges
//     result = edges?.reduce((acc: any, item: any) => {
//         const edgeValue = graphQlData[item.name];
//         if (!edgeValue) {
//             return acc;
//         }

//         acc[item.name] = {
//             label: edgeValue.name,
//             valueType: 'connect',
//             id: edgeValue.id
//         };

//         return acc;
//     }, result);

//     return result;

// };


// export const ApiDataConvertor = {
//     convertFormApi: convertDataFromGraphQl,
//     convertToApi: convertDataToGraphQl
// }