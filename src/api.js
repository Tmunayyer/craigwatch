import { spinnerState } from "./components/Spinner.vue";

export default async function (url, options) {
    spinnerState.setLoading(true);
    const response = await fetch(url, options);
    const body = await response.json();


    spinnerState.setLoading(false);
    return body;
}