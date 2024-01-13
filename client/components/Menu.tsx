import {
    DrawerContentComponentProps,
    DrawerContentScrollView,
} from "@react-navigation/drawer";
import { Image, Pressable, StyleSheet, Text, View } from "react-native";
import { useState } from "react";
import { useTheme } from "../context/theme";

const routes = ["Stores", "Documentation"];

export default function Menu(props: DrawerContentComponentProps) {
    const [index, setIndex] = useState(0);
    const { colors } = useTheme();

    function handleSelection(index: number) {
        setIndex(index);
        props.navigation.navigate(routes[index]);
    }

    return (
        <View style={[styles.container, { backgroundColor: colors.secondary }]}>
            <View style={styles.logoWrapper}>
                <Image
                    source={require("../assets/logo_transparent.png")}
                    style={styles.logo}
                />
                <Text style={styles.logoText}>Kache Krow</Text>
            </View>
            {routes.map((route, i) => (
                <Pressable
                    onPress={() => handleSelection(i)}
                    style={[
                        styles.menuItemWrapper,
                        index === i
                            ? {
                                  backgroundColor: colors.primary,
                                  ...styles.menuItemSelected,
                              }
                            : undefined,
                    ]}
                    disabled={index === i}
                    key={i}
                >
                    <Text
                        style={[
                            styles.menuItemText,
                            {
                                color:
                                    index === i ? colors.secondary : colors.text,
                            },
                        ]}
                    >
                        {route}
                    </Text>
                </Pressable>
            ))}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
    },
    logoWrapper: {
        flexDirection: "row",
        alignItems: "center",
    },
    logo: {
        width: 80,
        height: 80,
    },
    logoText: {
        fontFamily: "Alata, sans-serif",
        fontSize: 36,
        fontWeight: "600",
        fontStyle: "normal",
    },
    menuItemWrapper: {
        flexDirection: "row",
        maxWidth: "100%",
        height: 40,
        alignItems: "center",
        margin: 4,
    },
    menuItemSelected: {
        borderRadius: 8,
        color: "#FFFFFF",
    },
    menuItemText: {
        paddingLeft: 20,
        fontFamily: "Alata, sans-serif",
        fontSize: 18,
    },
});
