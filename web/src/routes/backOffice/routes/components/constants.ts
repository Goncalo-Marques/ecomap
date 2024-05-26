import { Style, Icon } from "ol/style";
import {
	CONTAINER_ICON_SRC,
	SELECTED_CONTAINER_ICON_SRC,
} from "../../../../lib/constants/map";

export const iconStyle = new Style({
	image: new Icon({
		src: CONTAINER_ICON_SRC,
	}),
});

export const selectedIconStyle = new Style({
	image: new Icon({
		src: SELECTED_CONTAINER_ICON_SRC,
	}),
});
