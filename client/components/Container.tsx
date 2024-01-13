import { PropsWithChildren } from "react";
import { Image, StyleSheet, View } from "react-native";
import { useTheme } from "../context/theme";

export default function ScreenContainer({ children }: PropsWithChildren) {
    const theme = useTheme();
    return (
        <View
            style={[
                styles.container,
                {
                    backgroundColor: "#000000"
                    
                },
            ]}
        >
            <Image
                source={require("../assets/background.png")}
                style={styles.background}
            />
            <View style={styles.contentWrapper}>
                {children}
            </View>
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        alignItems: "center",
        justifyContent: "center",
        padding: 20
    },
    contentWrapper: {
        maxWidth: 900,
        width: '95%',
        height: '100%'
    },
    background: {
        width: '100%',
        height: '100%',
        zIndex: -100,
        position: 'absolute',
        opacity: 0.2
    }
});
