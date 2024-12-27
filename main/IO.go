package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
	"strings"
)

// ReadCampusGraph 从 Redis 中读取数据并构建校园平面图
func ReadCampusGraph(redisClient *redis.Client) (*AdjList, error) {
	adjList := NewAdjList()

	nodeKeys := redisClient.Keys(context.Background(), "nodes:*").Val()
	for _, nodeKey := range nodeKeys {
		nodeIDStr := strings.TrimPrefix(nodeKey, "nodes:")
		nodeID, err := strconv.Atoi(nodeIDStr)
		if err != nil {
			return nil, fmt.Errorf("节点ID解析错误：%s", err)
		}

		nodeName, err := redisClient.Get(context.Background(), nodeKey).Result()
		if err != nil {
			return nil, fmt.Errorf("无法获取节点名称：%s", err)
		}

		adjList.Nodes[nodeID] = Node{
			ID:   nodeID,
			Name: nodeName,
		}
	}

	edgeKeys := redisClient.Keys(context.Background(), "edges:*").Val()
	for _, edgeKey := range edgeKeys {
		edgeIDStr := strings.TrimPrefix(edgeKey, "edges:")
		elements := strings.Split(edgeIDStr, "_")
		if len(elements) != 2 {
			return nil, fmt.Errorf("边ID解析错误：%s", edgeIDStr)
		}

		startID, err := strconv.Atoi(elements[0])
		if err != nil {
			return nil, fmt.Errorf("起始节点ID解析错误：%s", err)
		}
		endID, err := strconv.Atoi(elements[1])
		if err != nil {
			return nil, fmt.Errorf("结束节点ID解析错误：%s", err)
		}

		weightStr, err := redisClient.Get(context.Background(), edgeKey).Result()
		if err != nil {
			return nil, fmt.Errorf("无法获取边权重：%s", err)
		}
		weight, err := strconv.Atoi(weightStr)
		if err != nil {
			return nil, fmt.Errorf("边权重解析错误：%s", err)
		}

		edge := Edge{StartVex: startID, EndVex: endID, Weight: weight}
		adjList.AddEdge(edge)
	}

	return adjList, nil
}

// SaveCampusGraph 将校园平面图保存到 Redis 中
func SaveCampusGraph(adjList *AdjList, redisClient *redis.Client) error {
	// 清空当前数据库
	err := redisClient.FlushDB(context.Background()).Err()
	if err != nil {
		panic(err)
	}
	_ = redisClient.Set(context.Background(), "username", "noregret", 0).Err()
	_ = redisClient.Set(context.Background(), "admin_password", "Wang", 0).Err()

	// 保存节点信息
	for _, node := range adjList.Nodes {
		nodeKey := fmt.Sprintf("nodes:%d", node.ID)
		err := redisClient.Set(context.Background(), nodeKey, node.Name, 0).Err()
		if err != nil {
			return fmt.Errorf("保存节点信息失败：%s", err)
		}
	}

	// 保存边信息
	for _, node := range adjList.Nodes {
		edges := adjList.GetOutEdges(node.ID)
		for _, edge := range edges {
			edgeKey := fmt.Sprintf("edges:%d_%d", edge.StartVex, edge.EndVex)
			err := redisClient.Set(context.Background(), edgeKey, strconv.Itoa(edge.Weight), 0).Err()
			if err != nil {
				return fmt.Errorf("保存边信息失败：%s", err)
			}
		}
	}

	return nil
}
