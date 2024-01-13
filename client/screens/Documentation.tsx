import { StyleSheet, Text, View, ScrollView, TextInput, Image } from "react-native";
import ScreenContainer from "../components/Container";
import { useTheme } from "../context/theme";
import SyntaxHighlighter from "react-native-syntax-highlighter";
import { ThemeColors } from "../theme/colors";
import { Code } from "../components/Code";

export default function Documentation() {
    const { colors } = useTheme();

const code = `const example = () => {
    console.log(test)
}`;
    console.log(code)

    return (
        <ScreenContainer>
            <View style={[styles.header]}>
                <Text style={[styles.headerText, { color: colors.secondary }]}>
                    Documentation
                </Text>
                <Image source={require('../assets/paper.png')} style={styles.headerImage} resizeMode="contain"/>
            </View>
            <View style={[styles.container, styles.wrapper]}>
                <ScrollView
                    contentInsetAdjustmentBehavior="automatic"
                    style={[styles.scrollview]}
                >
                    <Code code={code} language={"javascript"}/>
                </ScrollView>
            </View>
        </ScreenContainer>
    );
}

const styles = StyleSheet.create({
    container: {
        width: "80%",
        padding: 20,
    },
    header: {
        flexDirection: 'row',
        width: "100%",
        alignItems: "center",
        padding: 20
    },
    headerImage: {
        width: 50,
        height: 50,
        marginLeft: 20
    },
    wrapper: {
        backgroundColor: ThemeColors.dark.secondary,
        padding: 15,
        borderRadius: 8,
        width: '100%'
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
