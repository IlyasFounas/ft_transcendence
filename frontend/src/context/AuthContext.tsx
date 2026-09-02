import { createContext, useContext, useState} from "react";
import type { ReactNode } from "react";

interface AuthContextType {
    accessToken: string | null;
    isAuthenticated: boolean;
    login: (accessToken: string) => void;
}

const AuthContext = createContext<AuthContextType>({
    accessToken: null,
    isAuthenticated: false,
    login: () => {},
});

function AuthProvider({children}: {children: ReactNode}) {
    
    const [accessToken, setAccessToken] = useState<string | null>(null);
    const isAuthenticated = accessToken !== null;
    const login = (accessToken: string) => {
        setAccessToken(accessToken);
    };

    return (
        <AuthContext.Provider 
            value={{
                accessToken,
                isAuthenticated,
                login,
            }}
        >
            {children}
        </AuthContext.Provider>
    );
}

function useAuth() {
    return useContext(AuthContext);
}

export { AuthProvider, useAuth };