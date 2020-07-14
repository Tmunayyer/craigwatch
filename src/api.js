export default async function (url, options) {
    const response = await fetch(url, options);
    const body = await response.json();

    return body;
}