// import { fullEntitiesMapVar, GetEntitiesNameCaption, nameCaptionEntitiesVar } from "@/hooks/use-entities";
// import { NameCaptionEntity } from "@/types/entity";
// import { AppStatusEvent, appStatusVar, entitiesInCodeGenerationVar } from "../app-state";
// import { ApolloClient } from "@apollo/client";
// import { addNotification } from "@/lib/notifications/store";
// import { AppStatusApplicationUpMessage } from "@/types/centrifugo";

// export const handleApplicationUp = (
//     message: AppStatusApplicationUpMessage,
//     client: ApolloClient<any>,
// ) => {

//     // clear full entities cache
//     //fullEntitiesMapVar({});

//     // invalidate apollo cache
//     client.getObservableQueries().forEach(query => {
//         client.cache.evict({ fieldName: query.queryName });
//         query.resetDiff();
//         query.resetLastResults();
//         query.reobserve();
//     });

//     client.refetchQueries({
//         include: "active",
//         optimistic: true
//     });

//     appStatusVar(AppStatusEvent.UP);

//     // Add notification
//     addNotification({
//         message: "Application is up and running",
//         type: 'info'
//     }, true);
// };