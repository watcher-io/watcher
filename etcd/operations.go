package etcd

import (
	"context"
	"github.com/watcher-io/watcher/model"
	"go.etcd.io/etcd/clientv3"
)

func FetchMember(
	ctx context.Context,
	c *clientv3.Client,
) (
	*model.Cluster,
	error,
) {
	var clusterState model.Cluster
	memberListResponse, err := c.MemberList(ctx)
	if err != nil {
		return nil, err
	}

	clusterState.ID = memberListResponse.Header.ClusterId
	for _, member := range memberListResponse.Members {
		statusResponse, err := c.Status(ctx, member.ClientURLs[0])
		if err != nil {
			return nil, err
		}
		clusterState.Leader = statusResponse.Leader
		clusterState.Members = append(clusterState.Members, model.ClusterMember{
			ID:         member.ID,
			Name:       member.Name,
			PeerURLS:   member.PeerURLs,
			ClientURLS: member.ClientURLs,
			IsLearner:  member.IsLearner,
			Status: model.MemberStatus{
				Version:          statusResponse.Version,
				DbSize:           statusResponse.DbSize,
				DbSizeInUse:      statusResponse.DbSizeInUse,
				RaftIndex:        statusResponse.RaftIndex,
				RaftTerm:         statusResponse.RaftTerm,
				RaftAppliedIndex: statusResponse.RaftAppliedIndex,
			},
		})
	}
	return &clusterState, err
}

func PutKV(
	ctx context.Context,
	c *clientv3.Client,
	putKVRequest *model.PutKVRequest,
) (
	*model.PutKVResponse,
	error,
) {

	putResponse, err := c.Put(ctx, putKVRequest.Key, putKVRequest.Value, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}
	kvResponse := model.PutKVResponse{
		MemberID:  putResponse.Header.GetMemberId(),
		RaftTerm:  putResponse.Header.GetRaftTerm(),
		ClusterID: putResponse.Header.GetClusterId(),
		NewKV: true,
	}
	if putResponse.PrevKv != nil {
		kvResponse.NewKV = false
		kvResponse.PreviousKV = model.PreviousKV{
			Key:            string(putResponse.PrevKv.Key),
			Version:        putResponse.PrevKv.Version,
			Value:          string(putResponse.PrevKv.Value),
			Lease:          putResponse.PrevKv.Lease,
		}
	}
	return &kvResponse, nil
}
