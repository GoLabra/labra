import { SnackbarProvider as NotistackProvider } from 'notistack';
import Defaults from "@/config/Defaults.json";

export const SnackbarProvider = ({ children }: { children: React.ReactNode }) => {

    return (<NotistackProvider maxSnack={Defaults.snackbar.maxSnack} anchorOrigin={{ vertical: 'top', horizontal: 'right' }}>
        {children}
    </NotistackProvider>)
}