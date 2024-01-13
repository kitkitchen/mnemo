import {
    Animated,
    StyleSheet,
    Text,
    TouchableOpacity,
    View,
} from "react-native";
import { CodeBlock, nord } from "react-code-blocks";
import * as Clipboard from "expo-clipboard";
import { ThemeColors } from "../theme/colors";
import { useRef } from "react";
import { Ionicons } from "@expo/vector-icons";

export const Code: React.FC<{ code: string; language: string }> = ({
    code,
    language,
}) => {
    const copyOpacity = useRef(new Animated.Value(0)).current;

    const copy = async () => {
        await Clipboard.setStringAsync(code);
        fade();
    };

    const fade = () => {
        Animated.timing(copyOpacity, {
            toValue: 100,
            duration: 2000,
            useNativeDriver: false,
        }).start();
        setTimeout(
            () =>
                Animated.timing(copyOpacity, {
                    toValue: 0,
                    duration: 2500,
                    useNativeDriver: false,
                }).start(),
            2000
        );
    };

    return (
        <View style={styles.container}>
            <View style={styles.button}>
                <Animated.View style={{ opacity: copyOpacity }}>
                    <Text style={[styles.buttonText]}>Copied!</Text>
                </Animated.View>
            <TouchableOpacity onPress={() => copy()}>
                <Ionicons name="copy-outline" size={28} color="#FFFFFF" />
            </TouchableOpacity>

            </View>
            <CodeBlock text={code} theme={nord} language={language} />
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        height: 200,
        position: "relative",
    },
    button: {
        position: "absolute",
        top: 5,
        right: 5,
        flexDirection: "row",
        alignItems: "center",
    },
    buttonText: {
        fontFamily: "Alata, sans-serif",
        color: ThemeColors.dark.secondary,
        marginRight: 10,
    },
});
