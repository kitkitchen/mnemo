import { createContext, useCallback, useContext, useState } from "react";
import { useThemeFonts, fonts } from "../theme/fonts";
import { ThemeColors } from "../theme/colors";

export enum ThemeType {
    dark = "dark",
    light = "light",
    spooky = "spooky"
}

type ThemeContextData = {
    themeType: ThemeType;
    colors: typeof ThemeColors.dark
    fonts: typeof fonts
};

export type ThemeContextState = ThemeContextData & {
    setThemeType: React.Dispatch<React.SetStateAction<ThemeType>>;
};

const ThemeContext = createContext<ThemeContextState | undefined>(undefined);

export function ThemeProvider({ children }: React.PropsWithChildren) {
    const [themeType, setThemeType] = useState<ThemeType>(ThemeType.dark);
    const fontsLoaded = useThemeFonts()
    console.log(fontsLoaded)
    
    const colors = themeType === ThemeType.dark ? ThemeColors.dark : ThemeColors.spooky

    const themeContextValue: ThemeContextState = {
        fonts,
        colors,
        themeType,
        setThemeType
    };

    return (
        <ThemeContext.Provider value={themeContextValue}>
            {children}
        </ThemeContext.Provider>
    );
}

export const useTheme = () => {
    const context = useContext(ThemeContext);

    if (!context) {
        throw new Error(
            "useTheme must be used within a ThemeContext provider"
        );
    }
    return context;
};
