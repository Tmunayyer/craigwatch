import { mount } from '@vue/test-utils';
import Metric from '../src/components/Metric.vue';

describe("should handle errors", () => {
    // ["label", "metricname", "data", "defaultSelected", "error"],
    const wrapper = mount(Metric, {
        propsData: {
            label: "bladerunner",
            metricname: "replicant",
            data: {},
            defaultSelect: "",
            error: true
        }
    });

    it("should render", () => {
        const container = wrapper.get(".metric-container");
        expect(container);

        const errorDiv = container.get(".error");
        expect(errorDiv);
    });
});
