package etcd

import (
	"context"
	"github.com/aka-achu/watcher/model"
	"go.etcd.io/etcd/clientv3"
)

func FetchMember(c *clientv3.Client) (model.Cluster, error){
	var clusterState model.Cluster
	memberListResponse, err := c.MemberList(context.Background())
	if err != nil {
		return clusterState, err
	}
	clusterState.ID = memberListResponse.Header.ClusterId
	for _, member := range memberListResponse.Members {
		statusResponse, err := c.Status(context.Background(), member.ClientURLs[0])
		if err != nil {
			return clusterState, err
		}
		clusterState.Leader = statusResponse.Leader
		clusterState.Members = append(clusterState.Members, model.ClusterMember{
			ID:         member.ID,
			Name:       member.Name,
			PeerURLS:   member.PeerURLs,
			ClientURLS: member.ClientURLs,
			IsLearner:  member.IsLearner,
			Status:     model.MemberStatus{
				Version:          statusResponse.Version,
				DbSize:           statusResponse.DbSize,
				DbSizeInUse:      statusResponse.DbSizeInUse,
				RaftIndex:        statusResponse.RaftIndex,
				RaftTerm:         statusResponse.RaftTerm,
				RaftAppliedIndex: statusResponse.RaftAppliedIndex,
			},
		})
	}
	return clusterState, err
}
