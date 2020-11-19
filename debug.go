package martinez_rueda

import "fmt"

func gatherSweepEventData(event SweepEvent) string {
	return fmt.Sprintf("index:%v is_left:%v x:%v y:%v other[x:%v y:%v]", event.id, event.isLeft, event.p.X(), event.p.Y(), event.other.p.X(), event.other.p.Y())
}

func gatherConnectorData(connector Connector) string {
	open_polygons := []string{}
	closed_polygons := []string{}

	for _, open := range connector.openPolygons {
		open_polygons = append(open_polygons, gatherPointChainData(*open))
	}

	for _, clo := range connector.closedPolygons {
		closed_polygons = append(closed_polygons, gatherPointChainData(*clo))
	}

	return fmt.Sprintf("closed:%v open_polygons:%v closed_polygons:%v", connector.isClosed(), open_polygons, closed_polygons)
}

func gatherPointChainData(chain PointChain) string {
	points := []string{}
	for _, seg := range chain.segments {
		points = append(points, fmt.Sprintf("x:%v y:%v", seg.X(), seg.Y()))
	}

	return fmt.Sprintf("closed:%v elements:%v", chain.closed, points)
}
