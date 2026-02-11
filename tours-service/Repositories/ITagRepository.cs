using TourService.Domain;

namespace TourService.Repositories;

public interface ITagRepository
{
    Task<List<Tag>> GetOrCreateAsync(IEnumerable<string> tagNames, CancellationToken ct);
}
