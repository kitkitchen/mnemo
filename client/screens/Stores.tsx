import { StyleSheet, Text, View, ScrollView, TextInput, Image } from "react-native";
import ScreenContainer from "../components/Container";
import { useTheme } from "../context/theme";
import SyntaxHighlighter from "react-native-syntax-highlighter";
import { ThemeColors } from "../theme/colors";
import { Code } from "../components/Code";

export default function Stores() {
    const { colors } = useTheme();

const code = `const example = () => {
    console.log(test)
}`;
    console.log(code)

    return (
        <ScreenContainer>
            <View style={[styles.header]}>
                <Text style={[styles.headerText, { color: colors.secondary }]}>
                    Stores
                </Text>
                <Image source={require('../assets/seeds.png')} style={styles.headerImage} resizeMode="contain"/>
            </View>
        </ScreenContainer>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        width: "100%",
    },
    header: {
        flexDirection: 'row',
        width: "100%",
        padding: 20,
        alignItems: "center",
    },
    headerImage: {
        marginLeft: 20,
        width: 50,
        height: 50
    },
    wrapper: {
        width: "100%",
        backgroundColor: ThemeColors.dark.secondary,
        padding: 15,
        marginBottom: 10,
        borderRadius: 8,
    },
    headerText: {
        fontFamily: "Alata, sans-serif",
        fontSize: 40,
    },
    scrollview: {
        width: "100%",
    },
    syntax: {
        backgroundColor: "#FFFFFF",
        margin: 10,
        shadowColor: "#000000",
        shadowOpacity: 0.4,
        shadowOffset: {
            width: 1,
            height: 1,
        },
        shadowRadius: 4,
    },
});
