import { AppStatus } from "@/lib/apollo/graphql.entities";
import { appStatusVar } from "../app-state";
import { addNotification } from "@/lib/notifications/store";
import { ApolloClient } from "@apollo/client";
import { fullEntitiesMapVar } from "@/hooks/use-entities";

export const handleApplicationStatusUP = (client: ApolloClient<any>) => {
    if(appStatusVar() === AppStatus.Up) {
        return;
    }

    appStatusVar(AppStatus.Up);

    // clear full entities cache
    fullEntitiesMapVar({});

    // invalidate apollo cache
    client.getObservableQueries().forEach(query => {
        client.cache.evict({ fieldName: query.queryName });
        query.resetDiff();
        query.resetLastResults();
        query.reobserve();
    });

    client.refetchQueries({
        include: "active",
        optimistic: true
    });
};

export const handleApplicationStatusGenerating = () => {
    if(appStatusVar() === AppStatus.Generating) {
        return;
    }

    appStatusVar(AppStatus.Generating);

    addNotification({
        message: "Please continue using the application. You'll be able to modify the entity once we're back into full color",
        type: 'info'
    }, true);
};

export const handleApplicationStatusReverting = () => {
    if(appStatusVar() === AppStatus.Reverting) {
        return;
    }

    appStatusVar(AppStatus.Reverting);

    addNotification({
        message: "Code generation is reverting",
        type: 'info'
    });
};

export const handleApplicationStatusRestarting = () => {
    if(appStatusVar() === AppStatus.Restarting) {
        return;
    }

    appStatusVar(AppStatus.Restarting);

    addNotification({
        message: "Application is shuting down/restarting",
        type: 'info',
    }, true);
};


export const handleApplicationStatusFailed = () => {
    if(appStatusVar() === AppStatus.Fatal) {
        return;
    }

    appStatusVar(AppStatus.Fatal);

    addNotification({
        message: "Code generation failed. All changes will be reverted",
        type: 'error'
    }, true);
};